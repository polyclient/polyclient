// Copyright (C) 2025 Juan Mesa and contributors
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3
// as published by the Free Software Foundation, with the Runtime
// Library Exception. See the COPYING.RUNTIME file for details.

package main

import (
	"context"
	"log"
	"os"

	"github.com/polyclient/polyclient/cmd/polyclient/commands"
	"github.com/polyclient/polyclient/internal/version"
	"github.com/urfave/cli/v3"
)

func main() {
	cmd := (&cli.Command{
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
	})

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
