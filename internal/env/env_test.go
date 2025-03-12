//

package env_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/polyclient/polyclient/internal/env"
	"github.com/polyclient/polyclient/internal/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEnvManager(t *testing.T) {
	t.Parallel()

	t.Run("singleton instance", func(t *testing.T) {
		t.Parallel()
		t.Parallel()

		manager1 := env.GetEnvManager()
		manager2 := env.GetEnvManager()

		assert.NotNil(t, manager1)
		assert.Same(t, manager1, manager2, "Expected the same instance from multiple calls")
	})
}

func TestEnvManagerGet(t *testing.T) {
	t.Parallel()

	t.Run("known environment variable", func(t *testing.T) {
		t.Parallel()

		testKey := env.EnvPolyClientPluginsDir
		expectedValue := "/test/plugins/dir"
		t.Setenv(testKey, expectedValue)

		manager := env.GetEnvManager()
		value, err := manager.Get(testKey)

		require.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	})

	t.Run("unknown environment variable", func(t *testing.T) {
		t.Parallel()

		manager := env.GetEnvManager()
		value, err := manager.Get("UNKNOWN_ENV_VAR")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "unknown environment variable")
		assert.Empty(t, value)
	})

	t.Run("unset environment variable", func(t *testing.T) {
		t.Parallel()

		testKey := env.EnvPolyClientSettingsFile
		t.Setenv(testKey, "")

		manager := env.GetEnvManager()
		value, err := manager.Get(testKey)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "is not set")
		assert.Empty(t, value)
	})
}

func TestEnvManagerSetup(t *testing.T) {
	t.Run("default values set", func(t *testing.T) {
		t.Parallel()
		t.Setenv(env.EnvPolyClientPluginsDir, "")
		t.Setenv(env.EnvPolyClientSettingsFile, "")
		t.Setenv(env.EnvPolyClientKeymapFile, "")

		manager := env.GetEnvManager()
		err := manager.Setup()

		require.NoError(t, err)

		pluginsDir := os.Getenv(env.EnvPolyClientPluginsDir)
		settingsFile := os.Getenv(env.EnvPolyClientSettingsFile)
		keymapFile := os.Getenv(env.EnvPolyClientKeymapFile)

		assert.NotEmpty(t, pluginsDir)
		assert.NotEmpty(t, settingsFile)
		assert.NotEmpty(t, keymapFile)

		pluginsDirInfo, err := os.Stat(pluginsDir)
		require.NoError(t, err)
		assert.True(t, pluginsDirInfo.IsDir())
	})

	t.Run("existing values preserved", func(t *testing.T) {
		t.Parallel()
		expectedPluginsDir := filepath.Join(t.TempDir(), "custom-plugins")
		expectedSettingsFile := "/custom/settings.json"
		expectedKeymapFile := "/custom/keymap.json"

		t.Setenv(env.EnvPolyClientPluginsDir, expectedPluginsDir)
		t.Setenv(env.EnvPolyClientSettingsFile, expectedSettingsFile)
		t.Setenv(env.EnvPolyClientKeymapFile, expectedKeymapFile)

		manager := env.GetEnvManager()
		err := manager.Setup()

		require.NoError(t, err)

		assert.Equal(t, expectedPluginsDir, os.Getenv(env.EnvPolyClientPluginsDir))
		assert.Equal(t, expectedSettingsFile, os.Getenv(env.EnvPolyClientSettingsFile))
		assert.Equal(t, expectedKeymapFile, os.Getenv(env.EnvPolyClientKeymapFile))

		pluginsDirInfo, err := os.Stat(expectedPluginsDir)
		require.NoError(t, err)
		assert.True(t, pluginsDirInfo.IsDir())
	})
}

func mockVersion(t *testing.T, mockValue string) func() {
	originalVersion := version.Version()

	t.Setenv("MOCK_VERSION", mockValue)

	return func() {
		t.Setenv("MOCK_VERSION", originalVersion)
	}
}

func TestDevMode(t *testing.T) {
	restore := mockVersion(t, "dev")
	defer restore()

	assert.Equal(t, "dev", version.Version())

	t.Setenv(env.EnvPolyClientPluginsDir, "")
	t.Setenv(env.EnvPolyClientSettingsFile, "")
	t.Setenv(env.EnvPolyClientKeymapFile, "")

	manager := env.GetEnvManager()
	err := manager.Setup()
	require.NoError(t, err)

	cwd, err := os.Getwd()
	require.NoError(t, err)

	expectedDevDir := filepath.Join(cwd, ".polyclientdev")

	pluginsDir := os.Getenv(env.EnvPolyClientPluginsDir)
	settingsFile := os.Getenv(env.EnvPolyClientSettingsFile)
	keymapFile := os.Getenv(env.EnvPolyClientKeymapFile)

	assert.Equal(t, filepath.Join(expectedDevDir, "plugins"), pluginsDir)
	assert.Equal(t, filepath.Join(expectedDevDir, "settings.json"), settingsFile)
	assert.Equal(t, filepath.Join(expectedDevDir, "keymap.json"), keymapFile)

	devDirInfo, err := os.Stat(expectedDevDir)
	require.NoError(t, err)
	assert.True(t, devDirInfo.IsDir())

	err = os.RemoveAll(expectedDevDir)
	require.NoError(t, err)
}

func TestProductionMode(t *testing.T) {
	restore := mockVersion(t, "1.0.0")
	defer restore()

	assert.Equal(t, "1.0.0", version.Version())

	t.Setenv(env.EnvPolyClientPluginsDir, "")
	t.Setenv(env.EnvPolyClientSettingsFile, "")
	t.Setenv(env.EnvPolyClientKeymapFile, "")

	manager := env.GetEnvManager()
	err := manager.Setup()
	require.NoError(t, err)

	userConfigDir, err := os.UserConfigDir()
	require.NoError(t, err)

	expectedConfigDir := filepath.Join(userConfigDir, "polyclient")

	pluginsDir := os.Getenv(env.EnvPolyClientPluginsDir)
	settingsFile := os.Getenv(env.EnvPolyClientSettingsFile)
	keymapFile := os.Getenv(env.EnvPolyClientKeymapFile)

	assert.Equal(t, filepath.Join(expectedConfigDir, "plugins"), pluginsDir)
	assert.Equal(t, filepath.Join(expectedConfigDir, "settings.json"), settingsFile)
	assert.Equal(t, filepath.Join(expectedConfigDir, "keymap.json"), keymapFile)

	configDirInfo, err := os.Stat(expectedConfigDir)
	require.NoError(t, err)
	assert.True(t, configDirInfo.IsDir())

	err = os.RemoveAll(expectedConfigDir)
	require.NoError(t, err)
}

func TestMain(m *testing.M) {
	exitCode := m.Run()

	userConfigDir, err := os.UserConfigDir()
	if err == nil {
		os.RemoveAll(filepath.Join(userConfigDir, "polyclient"))
	}

	cwd, err := os.Getwd()
	if err == nil {
		os.RemoveAll(filepath.Join(cwd, ".polyclientdev"))
	}

	os.Exit(exitCode)
}
