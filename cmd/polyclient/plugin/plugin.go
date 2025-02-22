// Copyright (C) 2025 Juan Mesa and contributors
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3
// as published by the Free Software Foundation, with the Runtime
// Library Exception. See the COPYING.RUNTIME file for details.

package plugin

import (
	"github.com/urfave/cli/v3"
)

func NewPluginCommand() *cli.Command {
	return &cli.Command{
		Name:     "plugin",
		Usage:    "Manage PolyClient plugins from the CLI",
		Category: "Plugins",
		Commands: []*cli.Command{
			newLoadCommand(),
		},
	}
}
