// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package env

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrepare(t *testing.T) {
	expectedConnections := filepath.Join(t.TempDir(), "connections")
	expectedPlugins := filepath.Join(t.TempDir(), "plugins")

	t.Setenv(VariableConnectionsDir.String(), expectedConnections)
	t.Setenv(VariablePluginsDir.String(), expectedPlugins)

	provider := NewSystemProvider()

	err := Prepare(provider)
	require.NoError(t, err)

	connectionsDir, err := provider.Get(VariableConnectionsDir)
	require.NoError(t, err)

	pluginsDir, err := provider.Get(VariablePluginsDir)
	require.NoError(t, err)

	expected := []string{
		expectedConnections,
		expectedPlugins,
	}

	actual := []string{
		connectionsDir,
		pluginsDir,
	}

	assert.ElementsMatch(t, expected, actual)
}
