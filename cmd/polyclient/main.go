// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

// Package main is the main entry point for the PolyClient CLI application.
package main

import (
	"context"
	"log"
	"os"

	pCli "github.com/polyclient/polyclient/cli"
	pCliDB "github.com/polyclient/polyclient/cli/db"
	"github.com/polyclient/polyclient/internal/engine"
	"github.com/polyclient/polyclient/internal/version"
	"github.com/urfave/cli/v3"
)

// main initializes and runs the PolyClient CLI application.
func main() {
	e, err := engine.NewEngine(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	cmd := &cli.Command{
		Name:                  "polyclient",
		Usage:                 "A command-line interface for PolyClient",
		Version:               version.Version(),
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			pCli.NewVersionCommand(e),
			pCli.NewDocsCommand(e),
			pCli.NewPluginCommand(e),
			pCli.NewGUICommand(e),
			pCli.NewLogCommand(e),
			pCliDB.NewDatabaseCommand(e),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
