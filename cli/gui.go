// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package cli

import (
	"context"

	"github.com/urfave/cli/v3"
)

// NewGUICommand creates a CLI command for launching PolyClient in GUI mode.
func NewGUICommand() *cli.Command {
	return &cli.Command{
		Name:  "gui",
		Usage: "Launch PolyClient in GUI mode",
		Action: func(context.Context, *cli.Command) error {
			return nil
		},
	}
}
