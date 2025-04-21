// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

// Package env provides environment configuration for PolyClient.
package env

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/polyclient/polyclient/internal/version"
)

const (
	// EnvPolyClientConnectionsDir controls the PolyClient connections directory.
	EnvPolyClientConnectionsDir = "POLYCLIENT_CONNECTIONS_DIR"

	// EnvPolyClientPluginsDir controls the PolyClient plugins directory.
	EnvPolyClientPluginsDir = "POLYCLIENT_PLUGINS_DIR"

	// EnvPolyClientSettingsFile controls the PolyClient settings file.
	EnvPolyClientSettingsFile = "POLYCLIENT_SETTINGS_FILE"

	// EnvPolyClientKeymapFile controls the PolyClient keymap file.
	EnvPolyClientKeymapFile = "POLYCLIENT_KEYMAP_FILE"
)

// Variable represents a PolyClient environment variable and its configuration.
type Variable struct {
	Name  string
	IsDir bool
}

// Manager manages environment variables and their initialization.
type Manager struct {
	vars map[string]Variable
}

var (
	globalManager *Manager
	once          sync.Once
)

// GetManager returns the singleton instance of Manager, ensuring it is initialized only once.
func GetManager() *Manager {
	once.Do(func() {
		globalManager = &Manager{
			vars: map[string]Variable{
				EnvPolyClientConnectionsDir: {EnvPolyClientConnectionsDir, true},
				EnvPolyClientPluginsDir:     {EnvPolyClientPluginsDir, true},
				EnvPolyClientSettingsFile:   {EnvPolyClientSettingsFile, false},
				EnvPolyClientKeymapFile:     {EnvPolyClientKeymapFile, false},
			},
		}
		_ = globalManager.Setup()
	})

	return globalManager
}

// Setup initializes environment variables with default values where necessary.
func (m *Manager) Setup() error {
	var errs []error

	for name, envVar := range m.vars {
		defaultPath, err := getDefaultPath(name)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get default path for %s: %w", name, err))

			continue
		}

		if err := ensureEnvVarSet(name, defaultPath); err != nil {
			errs = append(errs, fmt.Errorf("failed to setup %s: %w", name, err))
		}

		if envVar.IsDir {
			if err := ensureDirExists(defaultPath); err != nil {
				errs = append(
					errs,
					fmt.Errorf("failed to create directory %s: %w", defaultPath, err),
				)
			}
		}
	}

	return errors.Join(errs...)
}

// Get retrieves the value of a PolyClient environment variable.
func (m *Manager) Get(name string) (string, error) {
	if _, exists := m.vars[name]; !exists {
		return "", fmt.Errorf("unknown environment variable: %s", name)
	}

	path, err := getEnvPath(name)
	if err != nil {
		return "", fmt.Errorf("failed to get %s: %w", name, err)
	}

	if path == "" {
		return "", fmt.Errorf("environment variable %s is not set", name)
	}

	return path, nil
}

// getDefaultPath determines the default path based on the environment.
func getDefaultPath(name string) (string, error) {
	switch name {
	case EnvPolyClientConnectionsDir:
		return getUserConfigPath("connections")
	case EnvPolyClientPluginsDir:
		return getUserConfigPath("plugins")
	case EnvPolyClientSettingsFile:
		return getUserConfigPath("settings.json")
	case EnvPolyClientKeymapFile:
		return getUserConfigPath("keymap.json")
	default:
		return "", fmt.Errorf("unknown environment variable: %s", name)
	}
}

// getUserConfigPath returns the user config path, using a different path in dev mode.
func getUserConfigPath(subPath string) (string, error) {
	if subPath == "" {
		return "", errors.New("default path cannot be empty")
	}

	basePath, err := getConfigBasePath()
	if err != nil {
		return "", err
	}

	return filepath.Join(basePath, subPath), nil
}

// getConfigBasePath returns the base configuration directory.
// During development, it uses the .polyclientdev directory to simulate
// a production environment.
func getConfigBasePath() (string, error) {
	if version.Version() == "dev" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %w", err)
		}

		devDir := filepath.Join(cwd, ".polyclientdev")
		if err := os.MkdirAll(devDir, 0o750); err != nil {
			return "", fmt.Errorf("failed to create dev directory %s: %w", devDir, err)
		}

		return devDir, nil
	}

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}

	return filepath.Join(userConfigDir, "polyclient"), nil
}

// getEnvPath retrieves the value of an environment variable.
func getEnvPath(envVar string) (string, error) {
	path := os.Getenv(envVar)
	if path == "" {
		return "", fmt.Errorf("environment variable %s is not set", envVar)
	}

	return path, nil
}

// ensureEnvVarSet sets an environment variable if it is not already set.
func ensureEnvVarSet(envVar, defaultValue string) error {
	if _, exists := os.LookupEnv(envVar); exists {
		return nil
	}

	if err := os.Setenv(envVar, defaultValue); err != nil {
		return fmt.Errorf("failed to set %s: %w", envVar, err)
	}

	return nil
}

// ensureDirExists creates a directory if it does not exist.
func ensureDirExists(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	if err := os.MkdirAll(absPath, 0o750); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", absPath, err)
	}

	return nil
}
