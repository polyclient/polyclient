// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package cli

import (
	"context"

	"github.com/polyclient/polyclient/internal/application"
	"github.com/polyclient/polyclient/internal/plugin"
	"github.com/urfave/cli/v3"
)

// NewPluginCommand creates a CLI command for managing PolyClient plugins.
func NewPluginCommand(app *application.Application) *cli.Command {
	return &cli.Command{
		Name:  "plugin",
		Usage: "Manage PolyClient plugins from the CLI",
		Commands: []*cli.Command{
			newLoadCommand(app),
			newUnloadCommand(app),
		},
	}
}

// newLoadCommand creates a CLI command for loading PolyClient plugins.
func newLoadCommand(app *application.Application) *cli.Command {
	return &cli.Command{
		Name:  "load",
		Usage: "Load a plugin",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "path",
				Aliases:  []string{"p"},
				Usage:    "Path to the plugin directory or Wasm file",
				Required: true,
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			path := cmd.String("path")

			manifestPath, err := plugin.FindManifestPath(path)
			if err != nil {
				return err
			}

			if _, err := app.PluginsRegistry.LoadPlugin(manifestPath); err != nil {
				return err
			}

			return nil
		},
	}
}

// newUnloadCommand creates a CLI command for unloading PolyClient plugins.
func newUnloadCommand(app *application.Application) *cli.Command {
	return &cli.Command{
		Name:  "unload",
		Usage: "Unload a plugin",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "id",
				Usage:    "ID of the plugin to unload",
				Required: true,
			},
		},
		Action: func(_ context.Context, cmd *cli.Command) error {
			id := cmd.String("id")

			return app.PluginsRegistry.UnloadPlugin(id)
		},
	}
}
