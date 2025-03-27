// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package cli

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/polyclient/polyclient/api"
	"github.com/urfave/cli/v3"
)

// NewGUICommand creates a CLI command for launching the PolyClient GUI.
func NewGUICommand() *cli.Command {
	return &cli.Command{
		Name:  "gui",
		Usage: "Launch the PolyClient GUI",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			server := api.NewServer()
			url := "http://localhost:" + strconv.Itoa(server.Port)
			log.Printf("Starting server at: %s", url)

			sigCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
			defer stop()

			serverErr := make(chan error, 1)
			go func() {
				if err := server.ListenAndServe(); err != nil {
					serverErr <- err
				}
			}()

			if err := waitForServer(url); err != nil {
				return fmt.Errorf("failed to wait for server: %w", err)
			}

			fmt.Println("Opening browser at:", url)
			if err := openBrowser(url); err != nil {
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

func waitForServer(ur string) error {
	const maxAttempts = 10
	const delay = 100 * time.Millisecond

	for i := 0; i < maxAttempts; i++ {
		resp, err := http.Get(ur)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}

		time.Sleep(delay)
	}

	return fmt.Errorf("server not available after %d attempts", maxAttempts)
}

func openBrowser(url string) error {
	if runtime.GOOS == "windows" {
		return exec.Command("cmd", "/c", "start", url).Start()
	}

	if runtime.GOOS == "darwin" {
		return exec.Command("open", url).Start()
	}

	if isWSL() {
		return exec.Command("cmd", "/c", "start", url).Start()
	}

	return exec.Command("xdg-open", url).Start()
}
