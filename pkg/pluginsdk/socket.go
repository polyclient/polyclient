package pluginsdk

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetSocketPath returns the appropriate IPC path for the given plugin name.
// For Unix systems, it returns a Unix domain socket path.
// For Windows, it returns a named pipe path.
func GetSocketPath(pluginName string) string {
	switch runtime.GOOS {
	case "windows":
		return `\\.\pipe\` + pluginName
	default:
		return filepath.Join(os.TempDir(), pluginName+".sock")
	}
}
