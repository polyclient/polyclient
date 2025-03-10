// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

// Package env provides environment configuration for PolyClient.
package env

import (
	"fmt"
	"os"
	"path"

	"github.com/polyclient/polyclient/internal/version"
)

const (
	// EnvVarPolyClientDir is the environment variable for the PolyClient directory.
	EnvVarPolyClientDir = "POLYCLIENT_DIR"

	// EnvVarPolyClientPluginsDir is the environment variable for the PolyClient plugins directory.
	EnvVarPolyClientPluginsDir = "POLYCLIENT_PLUGINS_DIR"

	// EnvVarPolyClientSettingsFile is the environment variable for the PolyClient settings file.
	EnvVarPolyClientSettingsFile = "POLYCLIENT_SETTINGS_FILE"

	// EnvVarPolyClientKeymapFile is the environment variable for the PolyClient keymap file.
	EnvVarPolyClientKeymapFile = "POLYCLIENT_KEYMAP_FILE"
)

// Setup initializes the environment by creating config directories and setting env variables.
func Setup() error {
	if version.Version() == "dev" {
		return nil
	}

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to determine config directory: %w", err)
	}

	polyClientDir := path.Join(userConfigDir, ".polyclient")
	if err := os.Setenv(EnvVarPolyClientDir, polyClientDir); err != nil {
		return fmt.Errorf("failed to set PolyClient directory: %w", err)
	}
	if err := os.Setenv(EnvVarPolyClientPluginsDir, path.Join(polyClientDir, "plugins")); err != nil {
		return fmt.Errorf("failed to set plugins directory: %w", err)
	}
	if err := os.Setenv(EnvVarPolyClientSettingsFile, path.Join(polyClientDir, "settings.json")); err != nil {
		return fmt.Errorf("failed to set settings file path: %w", err)
	}
	if err := os.Setenv(EnvVarPolyClientKeymapFile, path.Join(polyClientDir, "keymap.json")); err != nil {
		return fmt.Errorf("failed to set keymap file path: %w", err)
	}

	const dirPerm = 0755

	if err := os.MkdirAll(polyClientDir, dirPerm); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.MkdirAll(path.Join(polyClientDir, "plugins"), dirPerm); err != nil {
		return fmt.Errorf("failed to create plugins directory: %w", err)
	}

	return nil
}
