// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package commands

import (
	"context"
	"fmt"
	"log"

	"github.com/polyclient/polyclient/internal/runtime/plugin"
	"github.com/urfave/cli/v3"
)

// NewPluginCommand creates a CLI command for managing PolyClient plugins.
func NewPluginCommand() *cli.Command {
	return &cli.Command{
		Name:  "plugin",
		Usage: "Manage PolyClient plugins from the CLI",
		Commands: []*cli.Command{
			newLoadCommand(),
		},
	}
}

func newLoadCommand() *cli.Command {
	return &cli.Command{
		Name:  "load",
		Usage: "Load a plugin",
		Action: func(context.Context, *cli.Command) error {
			lookupPaths := []string{
				"./plugins",
			}
			pr := plugin.NewPluginRegistry(lookupPaths)

			if err := pr.LoadPlugins(); err != nil {
				return err
			}

			plugin, err := pr.GetWASMPlugin("sqlite")
			if err != nil {
				return err
			}

			_, result, err := plugin.Call("greet", []byte("Juan"))
			if err != nil {
				return fmt.Errorf("failed to call function: %w", err)
			}

			log.Println(string(result))
			return nil
		},
	}
}
