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
	"github.com/polyclient/polyclient/drivers/postgres"
	"github.com/polyclient/polyclient/drivers/sqlite"
	"github.com/polyclient/polyclient/internal/db"
	"github.com/polyclient/polyclient/internal/env"
	"github.com/polyclient/polyclient/internal/plugin"
	"github.com/polyclient/polyclient/internal/sdk"
	"github.com/polyclient/polyclient/internal/version"
	"github.com/urfave/cli/v3"
)

// main initializes and runs the PolyClient CLI application.
func main() {
	if err := env.GetManager().Setup(); err != nil {
		log.Fatal("Error setting up environment:", err)
	}

	driversRegistry, err := loadDrivers()
	if err != nil {
		log.Fatal("Error loading database drivers:", err)
	}

	dbSDK := sdk.NewDatabaseSDK(driversRegistry)

	// pluginsRegistry, err := loadPlugins()
	// if err != nil {
	// 	log.Fatal("Error loading plugins:", err)
	// }

	cmd := &cli.Command{
		Name:                  "polyclient",
		Usage:                 "A command-line interface for PolyClient",
		Version:               version.Version(),
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			pCli.NewVersionCommand(),
			pCli.NewDocsCommand(),
			pCli.NewDatabaseCommand(dbSDK),
			// pCli.NewPluginCommand(pluginsRegistry),
			pCli.NewGUICommand(),
			pCli.NewLogCommand(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func loadDrivers() (*db.Registry[db.Driver], error) {
	dr := db.NewRegistry[db.Driver]()

	if err := dr.Register(sqlite.NewDriver()); err != nil {
		return nil, fmt.Errorf("failed to register SQLite driver: %w", err)
	}

	if err := dr.Register(postgres.NewDriver()); err != nil {
		return nil, fmt.Errorf("failed to register PostgreSQL driver: %w", err)
	}

	return dr, nil
}

// loadPlugins loads the built-in PolyClient plugins into the plugin registry.
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
