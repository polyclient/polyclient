// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

//nolint:testpackage  // Using same package name to be able to reset singleton manager
package env

import (
	"os"
	"path"
	"path/filepath"
	"sync"
	"testing"

	"github.com/polyclient/polyclient/internal/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetManagerForTest() {
	globalManager = nil
	once = sync.Once{}
}

func TestManagerSingleton(t *testing.T) {
	resetManagerForTest()

	manager1 := GetManager()
	manager2 := GetManager()

	assert.Same(t, manager1, manager2, "GetManager should return the same instance")
}

func TestManagerGet(t *testing.T) {
	resetManagerForTest()

	manager := GetManager()

	t.Run("Valid environment variables", func(t *testing.T) {
		variables := []string{
			EnvPolyClientPluginsDir,
			EnvPolyClientSettingsFile,
			EnvPolyClientKeymapFile,
		}

		for _, varName := range variables {
			value, err := manager.Get(varName)
			require.NoError(t, err, "Getting %s should not error", varName)
			assert.NotEmpty(t, value, "Value for %s should not be empty", varName)
		}
	})

	t.Run("Unknown environment variable", func(t *testing.T) {
		_, err := manager.Get("UNKNOWN_VAR")
		require.Error(t, err, "Getting unknown var should return error")
		assert.Contains(
			t,
			err.Error(),
			"unknown environment variable",
			"Error message should be descriptive",
		)
	})
}

func TestCustomEnvironmentValues(t *testing.T) {
	origPluginsDir := os.Getenv(EnvPolyClientPluginsDir)
	origSettingsFile := os.Getenv(EnvPolyClientSettingsFile)
	origKeymapFile := os.Getenv(EnvPolyClientKeymapFile)

	defer func() {
		t.Setenv(EnvPolyClientPluginsDir, origPluginsDir)
		t.Setenv(EnvPolyClientSettingsFile, origSettingsFile)
		t.Setenv(EnvPolyClientKeymapFile, origKeymapFile)
	}()

	customPluginsDir := path.Join(t.TempDir(), "plugins")
	customSettingsFile := path.Join(t.TempDir(), "settings.json")
	customKeymapFile := path.Join(t.TempDir(), "keymap.json")

	t.Setenv(EnvPolyClientPluginsDir, customPluginsDir)
	t.Setenv(EnvPolyClientSettingsFile, customSettingsFile)
	t.Setenv(EnvPolyClientKeymapFile, customKeymapFile)

	resetManagerForTest()

	manager := GetManager()

	pluginsDir, err := manager.Get(EnvPolyClientPluginsDir)
	require.NoError(t, err)
	assert.Equal(t, customPluginsDir, pluginsDir)

	settingsFile, err := manager.Get(EnvPolyClientSettingsFile)
	require.NoError(t, err)
	assert.Equal(t, customSettingsFile, settingsFile)

	keymapFile, err := manager.Get(EnvPolyClientKeymapFile)
	require.NoError(t, err)
	assert.Equal(t, customKeymapFile, keymapFile)
}

func TestDevModeConfigPath(t *testing.T) {
	origVersionFunc := version.Version
	origVersion := origVersionFunc()

	defer func() {
		version.SetVersion(origVersion)
	}()

	version.SetVersion("dev")

	_ = os.Unsetenv(EnvPolyClientPluginsDir)
	_ = os.Unsetenv(EnvPolyClientSettingsFile)
	_ = os.Unsetenv(EnvPolyClientKeymapFile)

	resetManagerForTest()

	manager := GetManager()

	cwd, err := os.Getwd()
	require.NoError(t, err)

	expectedDevDir := filepath.Join(cwd, ".polyclientdev")
	expectedPluginsDir := filepath.Join(expectedDevDir, "plugins")
	expectedSettingsFile := filepath.Join(expectedDevDir, "settings.json")
	expectedKeymapFile := filepath.Join(expectedDevDir, "keymap.json")

	pluginsDir, err := manager.Get(EnvPolyClientPluginsDir)
	require.NoError(t, err)
	assert.Equal(t, expectedPluginsDir, pluginsDir)

	settingsFile, err := manager.Get(EnvPolyClientSettingsFile)
	require.NoError(t, err)
	assert.Equal(t, expectedSettingsFile, settingsFile)

	keymapFile, err := manager.Get(EnvPolyClientKeymapFile)
	require.NoError(t, err)
	assert.Equal(t, expectedKeymapFile, keymapFile)

	_, err = os.Stat(expectedDevDir)
	require.NoError(t, err, "Dev directory should be created")

	_ = os.RemoveAll(expectedDevDir)
}

func TestProdModeConfigPath(t *testing.T) {
	origVersionFunc := version.Version
	origVersion := origVersionFunc()

	defer func() {
		version.SetVersion(origVersion)
	}()

	version.SetVersion("1.0.0")

	_ = os.Unsetenv(EnvPolyClientPluginsDir)
	_ = os.Unsetenv(EnvPolyClientSettingsFile)
	_ = os.Unsetenv(EnvPolyClientKeymapFile)

	resetManagerForTest()

	manager := GetManager()

	userConfigDir, err := os.UserConfigDir()
	require.NoError(t, err)

	expectedPluginsDir := filepath.Join(userConfigDir, "polyclient", "plugins")
	expectedSettingsFile := filepath.Join(userConfigDir, "polyclient", "settings.json")
	expectedKeymapFile := filepath.Join(userConfigDir, "polyclient", "keymap.json")

	pluginsDir, err := manager.Get(EnvPolyClientPluginsDir)
	require.NoError(t, err)
	assert.Equal(t, expectedPluginsDir, pluginsDir)

	settingsFile, err := manager.Get(EnvPolyClientSettingsFile)
	require.NoError(t, err)
	assert.Equal(t, expectedSettingsFile, settingsFile)

	keymapFile, err := manager.Get(EnvPolyClientKeymapFile)
	require.NoError(t, err)
	assert.Equal(t, expectedKeymapFile, keymapFile)

	_, err = os.Stat(expectedPluginsDir)
	require.NoError(t, err, "Config directory should be created")

	_ = os.RemoveAll(expectedPluginsDir)
}

func TestSetupEnsuresDirExists(t *testing.T) {
	_ = os.Unsetenv(EnvPolyClientPluginsDir)

	resetManagerForTest()

	manager := GetManager()

	pluginsDir, err := manager.Get(EnvPolyClientPluginsDir)
	require.NoError(t, err)

	info, err := os.Stat(pluginsDir)
	require.NoError(t, err, "Plugins directory should exist")
	assert.True(t, info.IsDir(), "Plugins path should be a directory")
}

func TestVariableDefinitions(t *testing.T) {
	resetManagerForTest()

	manager := GetManager()

	assert.Contains(t, manager.vars, EnvPolyClientPluginsDir)
	assert.Contains(t, manager.vars, EnvPolyClientSettingsFile)
	assert.Contains(t, manager.vars, EnvPolyClientKeymapFile)

	assert.True(t, manager.vars[EnvPolyClientPluginsDir].IsDir)
	assert.False(t, manager.vars[EnvPolyClientSettingsFile].IsDir)
	assert.False(t, manager.vars[EnvPolyClientKeymapFile].IsDir)
}

func TestEnvVarFallbacks(t *testing.T) {
	_ = os.Unsetenv(EnvPolyClientPluginsDir)
	_ = os.Unsetenv(EnvPolyClientSettingsFile)
	_ = os.Unsetenv(EnvPolyClientKeymapFile)

	resetManagerForTest()

	manager := GetManager()

	pluginsDir, err := manager.Get(EnvPolyClientPluginsDir)
	require.NoError(t, err)
	assert.NotEmpty(t, pluginsDir)

	settingsFile, err := manager.Get(EnvPolyClientSettingsFile)
	require.NoError(t, err)
	assert.NotEmpty(t, settingsFile)

	keymapFile, err := manager.Get(EnvPolyClientKeymapFile)
	require.NoError(t, err)
	assert.NotEmpty(t, keymapFile)

	assert.Equal(t, pluginsDir, os.Getenv(EnvPolyClientPluginsDir))
	assert.Equal(t, settingsFile, os.Getenv(EnvPolyClientSettingsFile))
	assert.Equal(t, keymapFile, os.Getenv(EnvPolyClientKeymapFile))
}
