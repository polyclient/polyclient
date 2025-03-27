// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package cli

import (
	"context"
	"log"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/polyclient/polyclient/api"
	"github.com/urfave/cli/v3"
)

// NewGUICommand creates a CLI command for launching PolyClient in GUI mode.
func NewGUICommand() *cli.Command {
	return &cli.Command{
		Name:  "gui",
		Usage: "Launch PolyClient in GUI mode",
		Action: func(context.Context, *cli.Command) error {
			server := api.NewServer()
			var wg sync.WaitGroup
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := server.Run(ctx); err != nil {
					log.Printf("Server error: %v", err)
				}
			}()

			go func() {
				time.Sleep(500 * time.Millisecond)
				if err := openBrowser("http://localhost:8081"); err != nil {
					log.Printf("Failed to open browser: %v", err)
				}
			}()

			<-ctx.Done()
			wg.Wait()

			return nil
		},
	}
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
