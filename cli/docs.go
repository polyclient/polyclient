// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package cli

import (
	"context"
	"os/exec"
	"runtime"

	"github.com/polyclient/polyclient/internal/application"
	"github.com/urfave/cli/v3"
)

// NewDocsCommand creates a CLI command for opening the documentation website.
func NewDocsCommand(app *application.Application) *cli.Command {
	return &cli.Command{
		Name:  "docs",
		Usage: "Open documentation website",
		Action: func(context.Context, *cli.Command) error {
			url := "https://polyclient.pages.dev"

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
		},
	}
}
