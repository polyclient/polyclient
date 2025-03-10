// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/polyclient/polyclient/cmd/polyclient/shell"
	"github.com/polyclient/polyclient/internal/env"
	"github.com/polyclient/polyclient/internal/runtime/plugin"
	"github.com/polyclient/polyclient/internal/version"
	"github.com/urfave/cli/v3"
)

// main initializes and runs the PolyClient CLI application.
func main() {
	if err := setupUserEnvironment(); err != nil {
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
			shell.NewVersionCommand(),
			shell.NewDatabaseCommand(pr),
			shell.NewPluginCommand(pr),
			shell.NewGuiCommand(),
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// setupUserEnvironment sets up the environment for PolyClient by creating the
// necessary config directories and setting the environment variables.
func setupUserEnvironment() error {
	if version.Version() == "dev" {
		return nil
	}

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to determine config directory: %w", err)
	}

	polyClientDir := path.Join(userConfigDir, ".polyclient")
	os.Setenv(env.EnvVarPolyClientDir, polyClientDir)
	os.Setenv(env.EnvVarPolyClientPluginsDir, path.Join(polyClientDir, "plugins"))
	os.Setenv(env.EnvVarPolyClientSettingsFile, path.Join(polyClientDir, "settings.json"))
	os.Setenv(env.EnvVarPolyClientKeymapFile, path.Join(polyClientDir, "keymap.json"))

	const dirPerm = 0755

	if err := os.MkdirAll(polyClientDir, dirPerm); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.MkdirAll(path.Join(polyClientDir, "plugins"), dirPerm); err != nil {
		return fmt.Errorf("failed to create plugins directory: %w", err)
	}

	return nil
}

// loadBuiltinPlugins loads the built-in PolyClient plugins into the plugin registry.
// The built-in plugins are loaded from the POLYCLIENT_PLUGINS_DIR environment variable.
func loadPlugins() (*plugin.PluginRegistry, error) {
	lookupPaths := []string{}

	if version.Version() == "dev" {
		lookupPaths = append(lookupPaths, "plugins")
	} else {
		lookupPaths = append(lookupPaths, os.Getenv(env.EnvVarPolyClientPluginsDir))
	}

	pr := plugin.NewPluginRegistry(lookupPaths)

	if err := pr.LoadPlugins(); err != nil {
		return nil, fmt.Errorf("failed to load plugins: %w", err)
	}

	return pr, nil
}
