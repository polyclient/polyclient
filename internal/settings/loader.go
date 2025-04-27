package settings

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"dario.cat/mergo"
	"github.com/polyclient/polyclient/internal/env"
)

//go:embed default_settings.json
var defaultSettingsFile embed.FS

// LoadFromEnv loads the user-defined settings from environment variables and
// deep merges them with the default settings from the default_settings.json embedded file.
func LoadFromEnv(envProvider env.Provider) (*Settings, error) {
	defaultSettings, err := loadDefaultSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to load default settings: %w", err)
	}

	userSettings, err := loadUserSettings(envProvider)
	if errors.Is(err, os.ErrNotExist) {
		userSettings = &Settings{}
	} else if err != nil {
		return nil, fmt.Errorf("failed to load user settings: %w", err)
	}

	if err := mergo.Merge(userSettings, defaultSettings); err != nil {
		return nil, fmt.Errorf("failed to merge settings: %w", err)
	}

	if err := userSettings.Validate(); err != nil {
		return nil, err
	}

	return userSettings, nil
}

// loadDefaultSettings loads the default settings from the default_settings.json embed file.
func loadDefaultSettings() (*Settings, error) {
	defaultSettingsBytes, err := defaultSettingsFile.ReadFile("default_settings.json")
	if err != nil {
		return nil, fmt.Errorf("failed to load default settings: %w", err)
	}

	var defaultSettings Settings
	if err := json.Unmarshal(defaultSettingsBytes, &defaultSettings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal default settings: %w", err)
	}

	return &defaultSettings, nil
}

// loadUserSettings loads the user settings from the POLYCLIENT_SETTINGS_FILE environment variable.
func loadUserSettings(envProvider env.Provider) (*Settings, error) {
	userSettingsFile, err := envProvider.Get(env.VariableSettingsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to get PolyClient settings file: %w", err)
	}

	absPath, err := filepath.Abs(userSettingsFile)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	userSettingsBytes, err := os.ReadFile(filepath.Clean(absPath))
	if err != nil {
		return nil, fmt.Errorf("failed to read PolyClient settings file: %w", err)
	}

	var userSettings Settings
	if err := json.Unmarshal(userSettingsBytes, &userSettings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal PolyClient settings: %w", err)
	}

	return &userSettings, nil
}
