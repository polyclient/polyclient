// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package cli

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/polyclient/polyclient/api"
	"github.com/polyclient/polyclient/internal/engine"
	"github.com/polyclient/polyclient/internal/system"
	"github.com/urfave/cli/v3"
)

// NewGUICommand creates a CLI command for launching the PolyClient GUI.
func NewGUICommand(e *engine.Engine) *cli.Command {
	return &cli.Command{
		Name:  "gui",
		Usage: "Launch the PolyClient GUI",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "host",
				Usage: "Host to listen on",
				Value: "127.0.0.1",
				Validator: func(host string) error {
					_, err := url.Parse(host)
					if err != nil {
						return fmt.Errorf("invalid host: %w", err)
					}

					return nil
				},
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "Port to listen on",
				Value:   8080,
				Validator: func(port int64) error {
					if port < 1 || port > 65535 {
						return errors.New("port must be between 1 and 65535")
					}

					return nil
				},
			},
			&cli.BoolFlag{
				Name:  "headless",
				Usage: "Run the GUI in headless mode",
				Value: false,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			flagHost := cmd.String("host")
			flagPort := cmd.Int("port")
			flagHeadless := cmd.Bool("headless")

			if flagHeadless {
				e.Settings.GUI.Enabled = false
			}

			server, err := api.NewServer(ctx, e,
				api.WithHost(flagHost),
				api.WithPort(int(flagPort)),
			)
			if err != nil {
				return fmt.Errorf("failed to create server: %w", err)
			}

			parsedURL, err := url.ParseRequestURI(fmt.Sprintf("http://%s:%d", flagHost, flagPort))
			if err != nil {
				return fmt.Errorf("invalid GUI URL: %w", err)
			}
			guiURL := parsedURL.String()

			sigCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
			defer stop()

			serverErr := make(chan error)
			go func() {
				if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					serverErr <- err
				}
			}()

			if err := waitForServer(sigCtx, guiURL); err != nil {
				return fmt.Errorf("failed to wait for server: %w", err)
			}

			e.Logger.Info("Server is running", "url", guiURL)
			if !flagHeadless {
				if err := system.OpenBrowser(guiURL); err != nil {
					e.Logger.Error("Failed to open browser", "error", err, "url", guiURL)
					fmt.Fprintf(os.Stderr, "Please open your browser at: %s\n", guiURL)
				}
			}

			if strings.TrimSpace(flagHost) == "0.0.0.0" {
				localIP, err := system.GetLocalIP()
				if err != nil {
					e.Logger.Warn("Could not determine local IP address", "error", err)
				} else {
					networkURL := fmt.Sprintf("http://%s:%d", localIP, flagPort)
					e.Logger.Info("Network access URL", "url", networkURL)
				}
			}

			select {
			case <-sigCtx.Done():
				e.Logger.Info("Shutting down server")
				if err := server.Shutdown(sigCtx); err != nil {
					return fmt.Errorf("failed to shutdown server: %w", err)
				}
			case err := <-serverErr:
				return fmt.Errorf("server failed: %w", err)
			}

			return nil
		},
	}
}

// waitForServer polls the server at guiURL until it responds with HTTP 200 OK or the context is canceled.
// It returns an error if the server is not available after a timeout.
func waitForServer(ctx context.Context, guiURL string) error {
	const timeout = 2 * time.Second
	const delay = 100 * time.Millisecond

	client := &http.Client{
		Timeout: 500 * time.Millisecond,
	}

	parsedURL, err := url.ParseRequestURI(guiURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), http.NoBody)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(delay)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return nil
		}

		time.Sleep(delay)
	}

	return fmt.Errorf("server not available after %s", timeout)
}
