package env

import (
	"errors"
	"fmt"
	"log/slog"
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

	// VariableLogLevel controls the PolyClient log level.
	VariableLogLevel Variable = "POLYCLIENT_LOG_LEVEL"
)

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
		return getConfigBaseSubpath("connections")
	case VariablePluginsDir:
		return getConfigBaseSubpath("plugins")
	case VariableSettingsFile:
		return getConfigBaseSubpath("settings.json")
	case VariableKeymapFile:
		return getConfigBaseSubpath("keymap.json")
	case VariableLogLevel:
		return slog.LevelInfo.String(), nil
	default:
		return "", nil
	}
}

// getConfigBaseSubpath returns the base path joined with the given subpath.
func getConfigBaseSubpath(subpath string) (string, error) {
	if subpath == "" {
		return "", errors.New("invalid empty config subpath")
	}

	basePath, err := getConfigBasePath()
	if err != nil {
		return "", err
	}

	return filepath.Join(basePath, subpath), nil
}

// getConfigBasePath returns the base path for PolyClient configuration.
func getConfigBasePath() (string, error) {
	if version.IsProd() {
		return getProdConfigBasePath()
	}

	return getDevConfigBasePath()
}

// getProdConfigBasePath returns the base path for PolyClient production environments.
func getProdConfigBasePath() (string, error) {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}

	return filepath.Join(userConfigDir, PolyClientProdConfigDir), nil
}

// getDevConfigBasePath returns the base path for PolyClient development environments.
func getDevConfigBasePath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	return filepath.Join(cwd, PolyClientDevConfigDir), nil
}
