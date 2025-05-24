// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package cli

import (
	"context"

	"github.com/polyclient/polyclient/internal/constant"
	"github.com/polyclient/polyclient/internal/engine"
	"github.com/polyclient/polyclient/internal/webbrowser"
	"github.com/urfave/cli/v3"
)

// NewDocsCommand creates a CLI command for opening the documentation website.
func NewDocsCommand(e *engine.Engine) *cli.Command {
	return &cli.Command{
		Name:  "docs",
		Usage: "Open documentation website",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return webbrowser.OpenURL(ctx, constant.DocsURL)
		},
	}
}
