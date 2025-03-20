// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

// Package main is the main entry point for the PolyClient CLI application.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	pCli "github.com/polyclient/polyclient/cli"
	"github.com/polyclient/polyclient/internal/env"
	"github.com/polyclient/polyclient/internal/plugin"
	"github.com/polyclient/polyclient/internal/version"
	"github.com/urfave/cli/v3"
)

// main initializes and runs the PolyClient CLI application.
func main() {
	if err := env.GetManager().Setup(); err != nil {
		log.Fatal("Error setting up environment:", err)
	}

	pr, err := loadPlugins()
	if err != nil {
		log.Fatal("Error loading plugins:", err)
	}

	cmd := &cli.Command{
		Name:                  "PolyClient CLI",
		Usage:                 "Manage and query your databases with ease",
		Version:               version.Version(),
		EnableShellCompletion: true,
		HideHelpCommand:       true,
		Commands: []*cli.Command{
			pCli.NewVersionCommand(),
			pCli.NewDocsCommand(),
			pCli.NewDatabaseCommand(pr),
			pCli.NewPluginCommand(pr),
			pCli.NewGUICommand(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// loadBuiltinPlugins loads the built-in PolyClient plugins into the plugin registry.
// The built-in plugins are loaded from the POLYCLIENT_PLUGINS_DIR environment variable.
func loadPlugins() (*plugin.Registry, error) {
	pluginsDir, err := env.GetManager().Get(env.EnvPolyClientPluginsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get PolyClient plugins directory: %w", err)
	}

	lookupPaths := []string{pluginsDir}

	pr, err := plugin.NewPluginRegistry(lookupPaths)
	if err != nil {
		return nil, fmt.Errorf("failed to create plugin registry: %w", err)
	}

	if err := pr.LoadPlugins(); err != nil {
		return nil, fmt.Errorf("failed to load plugins: %w", err)
	}

	return pr, nil
}
