// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

//nolint:testpackage  // Using same package name to access package internals
package env

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/polyclient/polyclient/internal/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	env1, err := NewManager()
	require.NoError(t, err)

	env2, err := NewManager()
	require.NoError(t, err)

	assert.NotSame(t, env1, env2, "New should return distinct instances")
}

func TestGet(t *testing.T) {
	t.Parallel()

	env, err := NewManager()
	require.NoError(t, err)

	t.Run("Valid environment variables", func(t *testing.T) {
		t.Parallel()

		variables := []string{
			PolyClientPluginsDir,
			PolyClientSettingsFile,
			PolyClientKeymapFile,
		}

		for _, varName := range variables {
			value, err := env.Get(varName)
			require.NoError(t, err)
			assert.NotEmpty(t, value, "Value for %s should not be empty", varName)
		}
	})

	t.Run("Unknown environment variable", func(t *testing.T) {
		_, err := env.Get("UNKNOWN_VAR")
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
	origPluginsDir := os.Getenv(PolyClientPluginsDir)
	origSettingsFile := os.Getenv(PolyClientSettingsFile)
	origKeymapFile := os.Getenv(PolyClientKeymapFile)

	defer func() {
		t.Setenv(PolyClientPluginsDir, origPluginsDir)
		t.Setenv(PolyClientSettingsFile, origSettingsFile)
		t.Setenv(PolyClientKeymapFile, origKeymapFile)
	}()

	customPluginsDir := path.Join(t.TempDir(), "plugins")
	customSettingsFile := path.Join(t.TempDir(), "settings.json")
	customKeymapFile := path.Join(t.TempDir(), "keymap.json")

	t.Setenv(PolyClientPluginsDir, customPluginsDir)
	t.Setenv(PolyClientSettingsFile, customSettingsFile)
	t.Setenv(PolyClientKeymapFile, customKeymapFile)

	manager, err := NewManager()
	require.NoError(t, err)

	pluginsDir, err := manager.Get(PolyClientPluginsDir)
	require.NoError(t, err)
	assert.Equal(t, customPluginsDir, pluginsDir)

	settingsFile, err := manager.Get(PolyClientSettingsFile)
	require.NoError(t, err)
	assert.Equal(t, customSettingsFile, settingsFile)

	keymapFile, err := manager.Get(PolyClientKeymapFile)
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

	_ = os.Unsetenv(PolyClientPluginsDir)
	_ = os.Unsetenv(PolyClientSettingsFile)
	_ = os.Unsetenv(PolyClientKeymapFile)

	manager, err := NewManager()
	require.NoError(t, err)

	cwd, err := os.Getwd()
	require.NoError(t, err)

	expectedDevDir := filepath.Join(cwd, ".polyclientdev")
	expectedPluginsDir := filepath.Join(expectedDevDir, "plugins")
	expectedSettingsFile := filepath.Join(expectedDevDir, "settings.json")
	expectedKeymapFile := filepath.Join(expectedDevDir, "keymap.json")

	pluginsDir, err := manager.Get(PolyClientPluginsDir)
	require.NoError(t, err)
	assert.Equal(t, expectedPluginsDir, pluginsDir)

	settingsFile, err := manager.Get(PolyClientSettingsFile)
	require.NoError(t, err)
	assert.Equal(t, expectedSettingsFile, settingsFile)

	keymapFile, err := manager.Get(PolyClientKeymapFile)
	require.NoError(t, err)
	assert.Equal(t, expectedKeymapFile, keymapFile)

	_, err = os.Stat(expectedDevDir)
	require.NoError(t, err)

	_ = os.RemoveAll(expectedDevDir)
}

func TestProdModeConfigPath(t *testing.T) {
	origVersionFunc := version.Version
	origVersion := origVersionFunc()

	defer func() {
		version.SetVersion(origVersion)
	}()

	version.SetVersion("1.0.0")

	_ = os.Unsetenv(PolyClientPluginsDir)
	_ = os.Unsetenv(PolyClientSettingsFile)
	_ = os.Unsetenv(PolyClientKeymapFile)

	manager, err := NewManager()
	require.NoError(t, err)

	userConfigDir, err := os.UserConfigDir()
	require.NoError(t, err)

	expectedPluginsDir := filepath.Join(userConfigDir, "polyclient", "plugins")
	expectedSettingsFile := filepath.Join(userConfigDir, "polyclient", "settings.json")
	expectedKeymapFile := filepath.Join(userConfigDir, "polyclient", "keymap.json")

	pluginsDir, err := manager.Get(PolyClientPluginsDir)
	require.NoError(t, err)
	assert.Equal(t, expectedPluginsDir, pluginsDir)

	settingsFile, err := manager.Get(PolyClientSettingsFile)
	require.NoError(t, err)
	assert.Equal(t, expectedSettingsFile, settingsFile)

	keymapFile, err := manager.Get(PolyClientKeymapFile)
	require.NoError(t, err)
	assert.Equal(t, expectedKeymapFile, keymapFile)

	_, err = os.Stat(expectedPluginsDir)
	require.NoError(t, err)

	_ = os.RemoveAll(expectedPluginsDir)
}

func TestSetupEnsuresDirExists(t *testing.T) {
	_ = os.Unsetenv(PolyClientPluginsDir)

	manager, err := NewManager()
	require.NoError(t, err)

	pluginsDir, err := manager.Get(PolyClientPluginsDir)
	require.NoError(t, err)

	info, err := os.Stat(pluginsDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir(), "Plugins path should be a directory")
}

func TestVariableDefinitions(t *testing.T) {
	manager, err := NewManager()
	require.NoError(t, err)

	assert.Contains(t, manager.vars, PolyClientPluginsDir)
	assert.Contains(t, manager.vars, PolyClientSettingsFile)
	assert.Contains(t, manager.vars, PolyClientKeymapFile)

	assert.True(t, manager.vars[PolyClientPluginsDir].IsDir)
	assert.False(t, manager.vars[PolyClientSettingsFile].IsDir)
	assert.False(t, manager.vars[PolyClientKeymapFile].IsDir)
}

func TestEnvVarFallbacks(t *testing.T) {
	_ = os.Unsetenv(PolyClientPluginsDir)
	_ = os.Unsetenv(PolyClientSettingsFile)
	_ = os.Unsetenv(PolyClientKeymapFile)

	manager, err := NewManager()
	require.NoError(t, err)

	pluginsDir, err := manager.Get(PolyClientPluginsDir)
	require.NoError(t, err)
	assert.NotEmpty(t, pluginsDir)

	settingsFile, err := manager.Get(PolyClientSettingsFile)
	require.NoError(t, err)
	assert.NotEmpty(t, settingsFile)

	keymapFile, err := manager.Get(PolyClientKeymapFile)
	require.NoError(t, err)
	assert.NotEmpty(t, keymapFile)

	assert.Equal(t, pluginsDir, os.Getenv(PolyClientPluginsDir))
	assert.Equal(t, settingsFile, os.Getenv(PolyClientSettingsFile))
	assert.Equal(t, keymapFile, os.Getenv(PolyClientKeymapFile))
}
