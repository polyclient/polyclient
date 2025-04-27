package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrepare(t *testing.T) {
	t.Setenv(string(VariableConnectionsDir), "/tmp/polyclient/connections")
	t.Setenv(string(VariablePluginsDir), "/tmp/polyclient/plugins")

	provider := NewSystemProvider()

	err := Prepare(provider)
	require.NoError(t, err)

	expected := []string{
		"/tmp/polyclient/connections",
		"/tmp/polyclient/plugins",
	}

	connectionsDir, err := provider.Get(VariableConnectionsDir)
	require.NoError(t, err)

	pluginsDir, err := provider.Get(VariablePluginsDir)
	require.NoError(t, err)

	actual := []string{
		connectionsDir,
		pluginsDir,
	}

	assert.ElementsMatch(t, expected, actual)
}
