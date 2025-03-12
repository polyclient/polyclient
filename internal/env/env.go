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
	// Environment variables used by PolyClient.
	EnvPolyClientPluginsDir   = "POLYCLIENT_PLUGINS_DIR"
	EnvPolyClientSettingsFile = "POLYCLIENT_SETTINGS_FILE"
	EnvPolyClientKeymapFile   = "POLYCLIENT_KEYMAP_FILE"
)

// EnvVar represents a PolyClient environment variable with associated functionality.
type EnvVar struct {
	Name  string
	IsDir bool
}

// EnvManager manages environment variables and their initialization.
type EnvManager struct {
	vars map[string]EnvVar
}

var (
	globalEnvManager *EnvManager
	once             sync.Once
)

// GetEnvManager returns the singleton instance of EnvManager, ensuring it is initialized only once.
func GetEnvManager() *EnvManager {
	once.Do(func() {
		globalEnvManager = &EnvManager{
			vars: map[string]EnvVar{
				EnvPolyClientPluginsDir:   {EnvPolyClientPluginsDir, true},
				EnvPolyClientSettingsFile: {EnvPolyClientSettingsFile, false},
				EnvPolyClientKeymapFile:   {EnvPolyClientKeymapFile, false},
			},
		}
		_ = globalEnvManager.Setup()
	})

	return globalEnvManager
}

// Setup initializes environment variables with default values where necessary.
func (m *EnvManager) Setup() error {
	var errs []error

	for name, envVar := range m.vars {
		defaultPath, err := getDefaultPath(name)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to get default path for %s: %w", name, err))
			continue
		}

		if err := ensureEnvVarSet(name, defaultPath, envVar.IsDir); err != nil {
			errs = append(errs, fmt.Errorf("failed to setup %s: %w", name, err))
		}
	}

	return errors.Join(errs...)
}

// Get retrieves the value of an environment variable.
func (m *EnvManager) Get(name string) (string, error) {
	if _, exists := m.vars[name]; !exists {
		return "", fmt.Errorf("unknown environment variable: %s", name)
	}

	return getEnvPath(name)
}

// getDefaultPath determines the default path based on the environment.
func getDefaultPath(name string) (string, error) {
	switch name {
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
func ensureEnvVarSet(envVar, defaultValue string, isDir bool) error {
	if os.Getenv(envVar) == "" {
		if err := os.Setenv(envVar, defaultValue); err != nil {
			return fmt.Errorf("failed to set %s: %w", envVar, err)
		}
	}

	absPath, err := filepath.Abs(defaultValue)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	if isDir {
		if err := os.MkdirAll(absPath, 0o750); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", absPath, err)
		}
	}

	return nil
}
