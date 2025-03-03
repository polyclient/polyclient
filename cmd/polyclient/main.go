// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package main

import (
	"context"
	"log"
	"os"

	"github.com/polyclient/polyclient/cmd/polyclient/commands"
	"github.com/polyclient/polyclient/internal/version"
	"github.com/urfave/cli/v3"
)

// main initializes and runs the PolyClient CLI application.
func main() {
	cmd := &cli.Command{
		Name:                  "PolyClient CLI",
		Usage:                 "Manage and query your databases with ease",
		Version:               version.Version(),
		EnableShellCompletion: true,
		HideHelpCommand:       true,
		Commands: []*cli.Command{
			commands.NewVersionCommand(),
			commands.NewDatabaseCommand(),
			commands.NewPluginCommand(),
			commands.NewGuiCommand(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
