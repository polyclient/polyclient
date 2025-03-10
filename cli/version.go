// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package cli

import (
	"context"
	"fmt"

	"github.com/polyclient/polyclient/internal/version"
	"github.com/urfave/cli/v3"
)

// NewVersionCommand creates a CLI command for displaying detailed version information.
func NewVersionCommand() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Show more detailed version information",
		Action: func(context.Context, *cli.Command) error {
			//nolint:forbidigo  // Printing version info to stdout is intentional here.
			fmt.Println(version.String())

			return nil
		},
	}
}
