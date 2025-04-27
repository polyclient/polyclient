package env

import (
	"fmt"
	"os"
)

// Prepare creates all the required directories and files for the PolyClient environment.
func Prepare(provider Provider) error {
	dirs := []string{}

	connectionsDir, err := provider.Get(VariableConnectionsDir)
	if err != nil {
		return fmt.Errorf("failed to determine connections directory: %w", err)
	}

	pluginsDir, err := provider.Get(VariablePluginsDir)
	if err != nil {
		return fmt.Errorf("failed to determine plugins directory: %w", err)
	}

	dirs = append(dirs, connectionsDir, pluginsDir)

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o750); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	return nil
}
