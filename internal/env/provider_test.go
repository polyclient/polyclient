// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package env

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/polyclient/polyclient/internal/version/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSystemProvider_Get(t *testing.T) {
	t.Run("returns env var when a value is set", func(t *testing.T) {
		expected := "/tmp/custom-connections-dir"

		t.Setenv(VariableConnectionsDir.String(), expected)

		provider := NewSystemProvider()

		value, err := provider.Get(VariableConnectionsDir)
		require.NoError(t, err)
		assert.Equal(t, expected, value)
	})

	t.Run("returns default value when env var is not set - dev mode", func(t *testing.T) {
		testutil.MockVersionDev(t)

		provider := NewSystemProvider()

		cwd, _ := os.Getwd()

		expectedConnectionsDir := filepath.Join(cwd, PolyClientDevConfigDir, "connections")
		connectionsDir, err := provider.Get(VariableConnectionsDir)
		require.NoError(t, err)
		assert.Equal(t, expectedConnectionsDir, connectionsDir)

		expectedPluginsDir := filepath.Join(cwd, PolyClientDevConfigDir, "plugins")
		pluginsDir, err := provider.Get(VariablePluginsDir)
		require.NoError(t, err)
		assert.Equal(t, expectedPluginsDir, pluginsDir)
	})

	t.Run("returns default value when env var is not set - prod mode", func(t *testing.T) {
		testutil.MockVersionProd(t)

		provider := NewSystemProvider()

		configDir, _ := os.UserConfigDir()

		expectedConnectionsDir := filepath.Join(configDir, PolyClientProdConfigDir, "connections")
		connectionsDir, err := provider.Get(VariableConnectionsDir)
		require.NoError(t, err)
		assert.Equal(t, expectedConnectionsDir, connectionsDir)

		expectedPluginsDir := filepath.Join(configDir, PolyClientProdConfigDir, "plugins")
		pluginsDir, err := provider.Get(VariablePluginsDir)
		require.NoError(t, err)
		assert.Equal(t, expectedPluginsDir, pluginsDir)
	})

	t.Run("returns empty string for unknown variable", func(t *testing.T) {
		provider := NewSystemProvider()

		value, err := provider.Get("UNKNOWN_VARIABLE")
		require.NoError(t, err)
		assert.Equal(t, "", value)
	})
}

func TestGetConfigPath(t *testing.T) {
	t.Run("returns error on empty subpath", func(t *testing.T) {
		_, err := getConfigPath("")
		require.Error(t, err)
	})

	t.Run("returns correct path in dev mode", func(t *testing.T) {
		testutil.MockVersionDev(t)

		cwd, _ := os.Getwd()

		expected := filepath.Join(cwd, PolyClientDevConfigDir, "connections")

		path, err := getConfigPath("connections")
		require.NoError(t, err)
		assert.Equal(t, expected, path)
	})

	t.Run("returns correct path in prod mode", func(t *testing.T) {
		testutil.MockVersionProd(t)

		configDir, _ := os.UserConfigDir()

		expected := filepath.Join(configDir, PolyClientProdConfigDir, "connections")

		path, err := getConfigPath("connections")
		require.NoError(t, err)
		assert.Equal(t, expected, path)
	})
}
