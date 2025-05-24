// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package env

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/polyclient/polyclient/internal/version"
)

// Variable represents an environment variable.
type Variable string

const (
	// VariableConnectionsDir controls the PolyClient connections directory.
	VariableConnectionsDir Variable = "POLYCLIENT_CONNECTIONS_DIR"

	// VariablePluginsDir controls the PolyClient plugins directory.
	VariablePluginsDir Variable = "POLYCLIENT_PLUGINS_DIR"

	// VariableSettingsFile controls the PolyClient settings file.
	VariableSettingsFile Variable = "POLYCLIENT_SETTINGS_FILE"

	// VariableKeymapFile controls the PolyClient keymap file.
	VariableKeymapFile Variable = "POLYCLIENT_KEYMAP_FILE"
)

// String returns the string representation of the variable.
func (v Variable) String() string {
	return string(v)
}

const (
	// PolyClientProdConfigDir controls the PolyClient production configuration directory.
	PolyClientProdConfigDir = "polyclient"

	// PolyClientDevConfigDir controls the PolyClient development configuration directory.
	PolyClientDevConfigDir = "test/testenv/polyclient"
)

// Provider defines environment variable access.
type Provider interface {
	Get(envVar Variable) (string, error)
}

// SystemProvider implements Provider using the real OS environment.
type SystemProvider struct{}

// NewSystemProvider creates and returns a new SystemProvider.
func NewSystemProvider() Provider {
	return &SystemProvider{}
}

// Get retrieves the environment variable value or fallback default.
func (p *SystemProvider) Get(envVar Variable) (string, error) {
	val := os.Getenv(string(envVar))
	if val != "" {
		return val, nil
	}

	defaultVal, err := resolveDefault(envVar)
	if err != nil {
		return "", err
	}

	return defaultVal, nil
}

// resolveDefault returns the default value for the given environment variable.
func resolveDefault(envVar Variable) (string, error) {
	switch envVar {
	case VariableConnectionsDir:
		return getConfigPath("connections")
	case VariablePluginsDir:
		return getConfigPath("plugins")
	case VariableSettingsFile:
		return getConfigPath("settings.json")
	case VariableKeymapFile:
		return getConfigPath("keymap.json")
	default:
		return "", nil
	}
}

// getConfigPath returns the configuration path for the given subpath.
func getConfigPath(subpath string) (string, error) {
	if subpath == "" {
		return "", errors.New("subpath cannot be empty")
	}

	if version.IsProd() {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user config directory: %w", err)
		}

		return filepath.Join(userConfigDir, PolyClientProdConfigDir, subpath), nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	return filepath.Join(cwd, PolyClientDevConfigDir, subpath), nil
}
