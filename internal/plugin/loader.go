package plugin

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/polyclient/polyclient/pkg/pluginsdk"
	"github.com/polyclient/polyclient/pkg/utils"
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

		manifests, err := pluginsdk.FindAllManifests(dir)
		if err != nil {
			// Log the error but continue processing other directories
			log.Printf("Error loading manifests for directory %s: %v\n", dir, err)
			continue
		}

		if len(manifests) == 0 {
			// Log the error but continue processing other directories
			log.Printf("No manifests found for directory %s\n", dir)
			continue
		}

		for _, manifestPath := range manifests {
			manifest, err := pluginsdk.LoadManifest(manifestPath)
			if err != nil {
				// Log the error but continue processing other directories
				log.Printf("Error loading manifest for directory %s: %v\n", dir, err)
				continue
			}

			executablePath := path.Join(
				strings.TrimSuffix(manifestPath, "manifest.json"),
				strings.TrimPrefix(manifest.EntryPoint, "./"),
			)
			if !isExecutable(executablePath) {
				// Log the error but continue processing other directories
				log.Printf("Invalid executable path for directory %s: %s\n", dir, executablePath)
				continue
			}

			if err := pm.LoadPlugin(manifest, executablePath); err != nil {
				// Log the error but continue processing other directories
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
func (pm *PluginManager) LoadPlugin(manifest *pluginsdk.Manifest, executablePath string) error {
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

	socketPath := utils.GetSocketPath(manifest.Name)

	attempts := 0
	maxAttempts := 3

	for {
		attempts++

		conn, err := createSocketClientConnection(socketPath)
		if err == nil {
			log.Printf("Plugin connected after %d attempts", attempts)

			client := pb.NewPluginClient(conn)

			pm.plugins[manifest.Name] = &managedPlugin{
				manifest: *manifest,
				client:   client,
				process:  cmd.Process,
				conn:     conn,
			}

			return nil
		}

		time.Sleep(time.Duration(attempts) * time.Second)

		if attempts == maxAttempts {
			cleanupSocketClientConnection(socketPath)
			cmd.Process.Kill()
			return fmt.Errorf("failed to create socket client connection after %d attempts: %w", attempts, err)
		}
	}
}

// ReloadPlugins reloads all available plugins. It is a wrapper around the LoadPlugins method.
func ReloadPlugins(pm *PluginManager) (int, error) {
	return pm.LoadPlugins()
}
