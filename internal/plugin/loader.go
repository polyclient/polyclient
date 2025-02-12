package plugin

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"strings"
	"time"

	"github.com/polyclient/polyclient/pkg/pluginsdk"
	pb "github.com/polyclient/polyclient/proto"
)

// LoadPlugins scans the plugin directories and loads all available plugins.
func (pm *PluginManager) LoadPlugins() (int, error) {
	// Clean up any existing plugins that might have been loaded
	if err := pm.cleanup(); err != nil {
		return 0, err
	}

	for _, dir := range pm.pluginsDirs {
		log.Print("Loading plugins from directory: " + dir)

		configPaths, err := pluginsdk.FindConfigPaths(dir)
		if err != nil || len(configPaths) == 0 {
			log.Printf("Failed to find plugin config files for directory %s: %v\n", dir, err)
			continue
		}

		for _, path := range configPaths {
			config, err := pluginsdk.LoadConfig(path)
			if err != nil {
				log.Printf("Failed to load plugin config for directory %s: %v\n", dir, err)
				continue
			}

			executablePath := filepath.Join(
				strings.TrimRight(path, "/"),
				strings.TrimPrefix(config.Manifest.EntryPoint, "./"),
			)
			if !isExecutable(executablePath) {
				log.Printf("Invalid executable path for directory %s: %s\n", dir, executablePath)
				continue
			}

			if err := pm.LoadPlugin(config, executablePath); err != nil {
				log.Printf("Error loading plugin for directory %s: %v\n", dir, err)
				continue
			}
		}
	}

	return len(pm.plugins), nil
}

// loadPlugin loads a plugin from the given path and registers it with the PluginManager.
// The plugin is started as a separate process and connected to via a socket file or named pipe.
// The method returns an error if the plugin process fails to start, the socket is not
// available, or the plugin fails to register.
func (pm *PluginManager) LoadPlugin(config *pluginsdk.Config, executablePath string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Start the plugin process
	cmd := exec.Command(executablePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start plugin process: %w", err)
	}

	log.Default().Printf("Plugin process started with PID %d", cmd.Process.Pid)

	socketPath := pluginsdk.GetSocketPath(config.Manifest.Name)

	attempts := 0
	maxAttempts := 3

	for {
		attempts++

		conn, err := createSocketClientConnection(socketPath)
		if err == nil {
			log.Printf("Plugin connected after %d attempts", attempts)

			client := pb.NewPluginClient(conn)

			pm.plugins[config.Manifest.Name] = &managedPlugin{
				config:  *config,
				client:  client,
				process: cmd.Process,
				conn:    conn,
			}

			return nil
		}

		time.Sleep(time.Duration(attempts) * time.Second)

		if attempts == maxAttempts {
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				log.Printf("failed to interrupt plugin process. killing process instead")
				log.Println(err)

				if err := cmd.Process.Kill(); err != nil {
					log.Printf("failed to kill plugin process")
					log.Println(err)
				}
			}

			// If the plugin fails to start, it may leave the socket file
			// behind. Interrupting the process should clean the socket file
			// but this code is a double check in case that fails.
			if err := cleanupSocketClientConnection(socketPath); err != nil {
				log.Printf("failed to cleanup socket client connection")
				log.Println(err)
			}

			return fmt.Errorf("failed to create socket client connection after %d attempts: %w", attempts, err)
		}
	}
}

// ReloadPlugins reloads all available plugins. It currently works as an alias for LoadPlugins
// but may be extended in the future to provide more fine-grained control over plugin reloading.
func ReloadPlugins(pm *PluginManager) (int, error) {
	return pm.LoadPlugins()
}
