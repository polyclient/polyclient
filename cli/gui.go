// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package cli

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/polyclient/polyclient/api"
	"github.com/polyclient/polyclient/internal/application"
	"github.com/urfave/cli/v3"
)

// NewGUICommand creates a CLI command for launching the PolyClient GUI.
func NewGUICommand(app *application.Application) *cli.Command {
	return &cli.Command{
		Name:  "gui",
		Usage: "Launch the PolyClient GUI",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			server, err := api.NewServer()
			if err != nil {
				return fmt.Errorf("failed to create server: %w", err)
			}

			guiURL := "http://localhost:" + strconv.Itoa(server.Port)
			log.Printf("Starting server at: %s", guiURL)

			sigCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
			defer stop()

			serverErr := make(chan error, 1)
			go func() {
				if err := server.ListenAndServe(); err != nil {
					serverErr <- err
				}
			}()

			if err := waitForServer(guiURL); err != nil {
				return fmt.Errorf("failed to wait for server: %w", err)
			}

			fmt.Println("Opening browser at:", guiURL)
			if err := openBrowser(guiURL); err != nil {
				log.Printf("Failed to open browser: %v", err)
			}

			fmt.Println("GUI launched successfully. Server is running. Press Ctrl+C to stop.")

			select {
			case <-sigCtx.Done():
				log.Println("Shutting down server...")
				if err := server.Shutdown(sigCtx); err != nil {
					log.Printf("Failed to shutdown server: %v", err)
				}

			case err := <-serverErr:
				return fmt.Errorf("server failed: %w", err)
			}

			return nil
		},
	}
}

func waitForServer(guiURL string) error {
	const maxAttempts = 10

	const delay = 100 * time.Millisecond

	for range maxAttempts {
		parsedURL, err := url.ParseRequestURI(guiURL)
		if err != nil {
			return fmt.Errorf("invalid URL: %w", err)
		}

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, parsedURL.String(), http.NoBody)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send request: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			if err := resp.Body.Close(); err != nil {
				log.Printf("Failed to close response body: %v", err)
			}

			return nil
		}

		time.Sleep(delay)
	}

	return fmt.Errorf("server not available after %d attempts", maxAttempts)
}

func openBrowser(guiURL string) error {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/c", "start", guiURL).Start()
	}

	if runtime.GOOS == "darwin" {
		return exec.Command("open", guiURL).Start()
	}

	if isWSL() {
		return exec.Command("cmd", "/c", "start", guiURL).Start()
	}

	return exec.Command("xdg-open", guiURL).Start()
}
