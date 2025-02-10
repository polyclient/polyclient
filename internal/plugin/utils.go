package plugin

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// cleanup closes all plugin connections and kills their processes.
func (pm *PluginManager) cleanup() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	var lastErr error

	for name, plugin := range pm.plugins {
		if err := plugin.conn.Close(); err != nil {
			lastErr = err
		} else {
			log.Default().Printf("[%s] gRPC connection closed", name)
		}

		if err := plugin.process.Kill(); err != nil {
			lastErr = err
		} else {
			log.Default().Printf("[%s] Process killed", name)
		}

	}

	pm.plugins = make(map[string]*managedPlugin)

	return lastErr
}

// isExecutable returns true if the given path is an executable file.
func isExecutable(path string) bool {
	if runtime.GOOS == "windows" {
		return filepath.Ext(path) == ".exe"
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.Mode()&0111 != 0
}
