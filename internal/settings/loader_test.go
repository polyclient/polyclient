package settings

import (
	"embed"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/polyclient/polyclient/internal/env"
)

//go:embed testdata/default_settings.json
var testdataFS embed.FS

type customEnvProvider struct {
	getFunc func(env.Variable) (string, error)
}

func (m *customEnvProvider) Get(key env.Variable) (string, error) {
	return m.getFunc(key)
}

func setupTempDir(t *testing.T) (string, func()) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "settings-test")

	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}

	return tempDir, func() { _ = os.RemoveAll(tempDir) }
}

func TestLoadFromEnvSuccess(t *testing.T) {
	origDefaultSettingsFile := defaultSettingsFile
	defaultSettingsFile = testdataFS
	defer func() { defaultSettingsFile = origDefaultSettingsFile }()

	tempDir, cleanup := setupTempDir(t)
	defer cleanup()

	userSettingsFile := filepath.Join(tempDir, "user_settings.json")
	userSettingsContent := `{
		"api": {
			"host": "localhost",
			"port": 8080
		},
		"logging": {
			"level": "debug"
		}
	}`
	if err := os.WriteFile(userSettingsFile, []byte(userSettingsContent), 0644); err != nil {
		t.Fatalf("failed to write user settings file: %v", err)
	}

	envProvider := &customEnvProvider{
		getFunc: func(v env.Variable) (string, error) {
			if v == env.VariableSettingsFile {
				return userSettingsFile, nil
			}
			return "", errors.New("unknown variable")
		},
	}

	settings, err := LoadFromEnv(envProvider)
	if err != nil {
		t.Errorf("LoadFromEnv failed: %v", err)
	}
	if settings == nil {
		t.Fatal("settings is nil")
	}

	if settings.API.Host != "localhost" {
		t.Errorf("expected API.Host to be %q, got %q", "localhost", settings.API.Host)
	}
	if settings.API.Port != 8080 {
		t.Errorf("expected API.Port to be %d, got %d", 8080, settings.API.Port)
	}
	if settings.Logging.Level != "debug" {
		t.Errorf("expected Logging.Level to be %q, got %q", "debug", settings.Logging.Level)
	}

	if !settings.Drivers.SQLite.Enabled {
		t.Error("expected Drivers.SQLite.Enabled to be true")
	}
	if settings.Logging.Format != "json" {
		t.Errorf("expected Logging.Format to be %q, got %q", "json", settings.Logging.Format)
	}
}

func TestLoadFromEnvNoUserSettings(t *testing.T) {
	origDefaultSettingsFile := defaultSettingsFile
	defaultSettingsFile = testdataFS
	defer func() { defaultSettingsFile = origDefaultSettingsFile }()

	tempDir, cleanup := setupTempDir(t)
	defer cleanup()

	nonExistentFile := filepath.Join(tempDir, "nonexistent.json")
	envProvider := &customEnvProvider{
		getFunc: func(v env.Variable) (string, error) {
			if v == env.VariableSettingsFile {
				return nonExistentFile, nil
			}
			return "", errors.New("unknown variable")
		},
	}

	settings, err := LoadFromEnv(envProvider)
	if err != nil {
		t.Errorf("LoadFromEnv failed: %v", err)
	}
	if settings == nil {
		t.Fatal("settings is nil")
	}

	if !settings.Drivers.SQLite.Enabled {
		t.Error("expected Drivers.SQLite.Enabled to be true")
	}
	if settings.Logging.Format != "json" {
		t.Errorf("expected Logging.Format to be %q, got %q", "json", settings.Logging.Format)
	}
	if settings.Logging.Level != "info" {
		t.Errorf("expected Logging.Level to be %q, got %q", "info", settings.Logging.Level)
	}
}

func TestLoadFromEnvInvalidUserSettings(t *testing.T) {
	origDefaultSettingsFile := defaultSettingsFile
	defaultSettingsFile = testdataFS
	defer func() { defaultSettingsFile = origDefaultSettingsFile }()

	tempDir, cleanup := setupTempDir(t)
	defer cleanup()

	userSettingsFile := filepath.Join(tempDir, "invalid_settings.json")
	if err := os.WriteFile(userSettingsFile, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("failed to write invalid user settings file: %v", err)
	}

	envProvider := &customEnvProvider{
		getFunc: func(v env.Variable) (string, error) {
			if v == env.VariableSettingsFile {
				return userSettingsFile, nil
			}
			return "", errors.New("unknown variable")
		},
	}

	settings, err := LoadFromEnv(envProvider)
	if err == nil {
		t.Error("expected LoadFromEnv to return an error")
	}
	if !strings.Contains(err.Error(), "failed to unmarshal PolyClient settings") {
		t.Errorf("expected error to contain %q, got %q", "failed to unmarshal PolyClient settings", err.Error())
	}
	if settings != nil {
		t.Error("expected settings to be nil")
	}
}

func TestLoadFromEnvInvalidPath(t *testing.T) {
	origDefaultSettingsFile := defaultSettingsFile
	defaultSettingsFile = testdataFS
	defer func() { defaultSettingsFile = origDefaultSettingsFile }()

	invalidPath := "/invalid/path/\x00/invalid.json"
	envProvider := &customEnvProvider{
		getFunc: func(v env.Variable) (string, error) {
			if v == env.VariableSettingsFile {
				return invalidPath, nil
			}
			return "", errors.New("unknown variable")
		},
	}

	settings, err := LoadFromEnv(envProvider)
	if err == nil {
		t.Error("expected LoadFromEnv to return an error")
	}
	if !strings.Contains(err.Error(), "invalid path") {
		t.Errorf("expected error to contain %q, got %q", "invalid path", err.Error())
	}
	if settings != nil {
		t.Error("expected settings to be nil")
	}
}

func TestLoadFromEnvEnvProviderError(t *testing.T) {
	origDefaultSettingsFile := defaultSettingsFile
	defaultSettingsFile = testdataFS
	defer func() { defaultSettingsFile = origDefaultSettingsFile }()

	envProvider := &customEnvProvider{
		getFunc: func(v env.Variable) (string, error) {
			return "", errors.New("env provider error")
		},
	}

	settings, err := LoadFromEnv(envProvider)
	if err == nil {
		t.Error("expected LoadFromEnv to return an error")
	}
	if !strings.Contains(err.Error(), "failed to get PolyClient settings file") {
		t.Errorf("expected error to contain %q, got %q", "failed to get PolyClient settings file", err.Error())
	}
	if settings != nil {
		t.Error("expected settings to be nil")
	}
}

func TestLoadDefaultSettingsSuccess(t *testing.T) {
	origDefaultSettingsFile := defaultSettingsFile
	defaultSettingsFile = testdataFS
	defer func() { defaultSettingsFile = origDefaultSettingsFile }()

	settings, err := loadDefaultSettings()
	if err != nil {
		t.Errorf("loadDefaultSettings failed: %v", err)
	}
	if settings == nil {
		t.Fatal("settings is nil")
	}

	if !settings.Drivers.SQLite.Enabled {
		t.Error("expected Drivers.SQLite.Enabled to be true")
	}
	if settings.Logging.Format != "json" {
		t.Errorf("expected Logging.Format to be %q, got %q", "json", settings.Logging.Format)
	}
	if settings.Logging.Level != "info" {
		t.Errorf("expected Logging.Level to be %q, got %q", "info", settings.Logging.Level)
	}
}

func TestLoadDefaultSettingsFileNotFound(t *testing.T) {
	origDefaultSettingsFile := defaultSettingsFile
	defaultSettingsFile = embed.FS{}
	defer func() { defaultSettingsFile = origDefaultSettingsFile }()

	settings, err := loadDefaultSettings()
	if err == nil {
		t.Error("expected loadDefaultSettings to return an error")
	}
	if !strings.Contains(err.Error(), "failed to load default settings") {
		t.Errorf("expected error to contain %q, got %q", "failed to load default settings", err.Error())
	}
	if settings != nil {
		t.Error("expected settings to be nil")
	}
}

func TestLoadDefaultSettingsInvalidJSON(t *testing.T) {
	origDefaultSettingsFile := defaultSettingsFile
	defaultSettingsFile = embed.FS{}
	defer func() { defaultSettingsFile = origDefaultSettingsFile }()

	settings, err := loadDefaultSettings()
	if err == nil {
		t.Error("expected loadDefaultSettings to return an error")
	}
	if !strings.Contains(err.Error(), "failed to load default settings") {
		t.Errorf("expected error to contain %q, got %q", "failed to load default settings", err.Error())
	}
	if settings != nil {
		t.Error("expected settings to be nil")
	}
}

func TestValidateSettings(t *testing.T) {
	origDefaultSettingsFile := defaultSettingsFile
	defaultSettingsFile = testdataFS
	defer func() { defaultSettingsFile = origDefaultSettingsFile }()

	settings, err := loadDefaultSettings()
	if err != nil {
		t.Fatalf("loadDefaultSettings failed: %v", err)
	}

	if err := settings.Validate(); err != nil {
		t.Errorf("Validate failed: %v", err)
	}

	invalidSettings := &Settings{
		API: API{
			Port: -1,
		},
	}
	if err := invalidSettings.Validate(); err == nil {
		t.Error("expected Validate to return an error")
	} else if !strings.Contains(err.Error(), "failed to validate settings") {
		t.Errorf("expected error to contain %q, got %q", "failed to validate settings", err.Error())
	}
}
