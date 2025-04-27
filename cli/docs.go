// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package cli

import (
	"context"

	"github.com/polyclient/polyclient/internal/engine"
	"github.com/polyclient/polyclient/internal/system"
	"github.com/urfave/cli/v3"
)

// NewDocsCommand creates a CLI command for opening the documentation website.
func NewDocsCommand(e *engine.Engine) *cli.Command {
	return &cli.Command{
		Name:  "docs",
		Usage: "Open documentation website",
		Action: func(context.Context, *cli.Command) error {
			const url = "https://polyclient.pages.dev"

			return system.OpenBrowser(url)
		},
	}
}
