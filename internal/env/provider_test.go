package env

import (
	"os"
	"testing"

	"log/slog"

	"github.com/polyclient/polyclient/internal/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemProvider_Get(t *testing.T) {
	t.Run("WhenVariableIsSet", func(t *testing.T) {
		provider := NewSystemProvider()
		key := string(VariableConnectionsDir)
		expected := "/tmp/polyclient/connections"

		t.Setenv(key, expected)

		val, err := provider.Get(VariableConnectionsDir)

		require.NoError(t, err)
		assert.Equal(t, expected, val)
	})

	t.Run("WithUnknownVariable", func(t *testing.T) {
		provider := NewSystemProvider()

		_ = os.Unsetenv("NON_EXISTENT_VAR")

		val, err := provider.Get(Variable("NON_EXISTENT_VAR"))

		require.NoError(t, err)
		assert.Empty(t, val)
	})
}

func TestSystemProvider_DefaultValues(t *testing.T) {
	t.Run("ConnectionsDir", func(t *testing.T) {
		provider := NewSystemProvider()
		_ = os.Unsetenv(string(VariableConnectionsDir))

		val, err := provider.Get(VariableConnectionsDir)

		require.NoError(t, err)
		assert.NotEmpty(t, val)
		assert.Contains(t, val, "connections")
	})

	t.Run("PluginsDir", func(t *testing.T) {
		provider := NewSystemProvider()
		_ = os.Unsetenv(string(VariablePluginsDir))

		val, err := provider.Get(VariablePluginsDir)

		require.NoError(t, err)
		assert.NotEmpty(t, val)
		assert.Contains(t, val, "plugins")
	})

	t.Run("SettingsFile", func(t *testing.T) {
		provider := NewSystemProvider()
		_ = os.Unsetenv(string(VariableSettingsFile))

		val, err := provider.Get(VariableSettingsFile)

		require.NoError(t, err)
		assert.NotEmpty(t, val)
		assert.Contains(t, val, "settings.json")
	})

	t.Run("KeymapFile", func(t *testing.T) {
		provider := NewSystemProvider()
		_ = os.Unsetenv(string(VariableKeymapFile))

		val, err := provider.Get(VariableKeymapFile)

		require.NoError(t, err)
		assert.NotEmpty(t, val)
		assert.Contains(t, val, "keymap.json")
	})

	t.Run("LogLevel", func(t *testing.T) {
		provider := NewSystemProvider()
		_ = os.Unsetenv(string(VariableLogLevel))

		val, err := provider.Get(VariableLogLevel)

		require.NoError(t, err)
		assert.Equal(t, slog.LevelInfo.String(), val)
	})
}

func TestConfigBasePath(t *testing.T) {
	t.Run("WhenVersionIsDev", func(t *testing.T) {
		oldVersion := version.Version()
		defer func() { version.SetVersion(oldVersion) }()

		version.SetVersion("dev")

		path, err := getConfigBasePath()

		require.NoError(t, err)
		assert.Contains(t, path, PolyClientDevConfigDir)
	})

	t.Run("WhenVersionIsProd", func(t *testing.T) {
		oldVersion := version.Version()
		defer func() { version.SetVersion(oldVersion) }()

		version.SetVersion("v1.2.3")

		path, err := getConfigBasePath()

		require.NoError(t, err)
		assert.Contains(t, path, PolyClientProdConfigDir)
	})
}

func TestConfigBaseSubpath(t *testing.T) {
	t.Run("WhenSubpathIsEmpty", func(t *testing.T) {
		path, err := getConfigBaseSubpath("")

		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid empty config subpath")
		assert.Empty(t, path)
	})
}
