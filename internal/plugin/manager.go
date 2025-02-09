package plugin

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/polyclient/polyclient/pkg/utils"
	pb "github.com/polyclient/polyclient/proto"
	"google.golang.org/grpc"
)

var (
	LoadingTimeout   time.Duration = 5 * time.Second
	ExecutionTimeout time.Duration = 30 * time.Second
)

// PluginManager is a manager for PolyClient plugins.
type PluginManager struct {
	pluginsDirs []string
	plugins     map[string]*managedPlugin
	mu          sync.RWMutex
}

// managedPlugin represents a plugin managed by the PluginManager and its associated resources.
type managedPlugin struct {
	info    *pb.PluginInfo
	client  pb.PluginClient
	process *os.Process
	conn    *grpc.ClientConn
}

// NewPluginManagerOptions defines the options for creating a new PluginManager.
type NewPluginManagerOptions struct {
	PluginsDirs []string
}

// NewPluginManager creates a new PluginManager instance.
func NewPluginManager(opts NewPluginManagerOptions) *PluginManager {
	return &PluginManager{
		pluginsDirs: opts.PluginsDirs,
		plugins:     make(map[string]*managedPlugin),
	}
}

// GetPlugin returns the plugin with the given name.
func (pm *PluginManager) GetPlugin(name string) (*managedPlugin, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	plugin, exists := pm.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return plugin, nil
}

// ListPlugins returns a list of all registered plugins and their information.
func (pm *PluginManager) GetPlugins() map[string]*managedPlugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	return pm.plugins
}

// LoadPlugins scans the plugin directories and loads all available plugins.
func (pm *PluginManager) LoadPlugins() (int, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if err := pm.cleanup(); err != nil {
		return 0, err
	}

	for _, dir := range pm.pluginsDirs {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() || !info.Mode().IsRegular() {
				return nil
			}

			if !isExecutable(path) {
				return nil
			}

			log.Default().Println("Loading plugin:", path)

			if err := pm.loadPlugin(path); err != nil {
				return fmt.Errorf("failed to load plugin %s: %w", path, err)
			}

			return nil
		})

		if err != nil {
			// Log the error but continue processing other directories
			fmt.Printf("Error walking directory %s: %v\n", dir, err)
		}
	}

	return len(pm.plugins), nil
}

// ExecuteAction executes the specified action on the plugin with the given name.
func (pm *PluginManager) Execute(pluginName, action string, payload []byte, metadata map[string]string) ([]byte, error) {
	pm.mu.Lock()
	plugin, exists := pm.plugins[pluginName]
	pm.mu.Unlock()

	if !exists {
		return nil, fmt.Errorf("plugin %s not found", pluginName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), ExecutionTimeout)
	defer cancel()

	result, err := plugin.client.Execute(ctx, &pb.PluginRequest{
		Action:   action,
		Payload:  payload,
		Metadata: metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute action %s: %w", action, err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf("plugin failed to execute action %s", result.Error)
	}

	return result.Payload, nil
}

// Shutdown stops all plugin processes and closes all plugin connections.
func (pm *PluginManager) Shutdown() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	return pm.cleanup()
}

// loadPlugin loads a plugin from the given path and registers it with the PluginManager.
// The plugin is started as a separate process and connected to via a Unix domain socket.
// The method returns an error if the plugin process fails to start, the socket is not
// available, or the plugin fails to register.
func (pm *PluginManager) loadPlugin(path string) error {
	// Start the plugin process
	cmd := exec.Command(path)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start plugin process: %w", err)
	}

	socketPath := utils.GetSocketPath(path)

	// Connect to the plugin
	conn, err := createSocketClient(socketPath)
	if err != nil {
		if err := cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill plugin process: %w", err)
		}

		return fmt.Errorf("failed to connect to plugin: %w", err)
	}

	client := pb.NewPluginClient(conn)

	// Register the plugin in the client
	ctx, cancel := context.WithTimeout(context.Background(), LoadingTimeout)
	defer cancel()

	info, err := client.Register(ctx, &pb.PluginInfo{})
	if err != nil {
		conn.Close()

		if err := cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill plugin process: %w", err)
		}

		return fmt.Errorf("failed to register plugin: %w", err)
	}

	pm.plugins[path] = &managedPlugin{
		info:    info,
		client:  client,
		process: cmd.Process,
		conn:    conn,
	}

	return nil
}

// cleanup closes all plugin connections and kills their processes.
func (pm *PluginManager) cleanup() error {
	var lastErr error

	for _, plugin := range pm.plugins {
		if err := plugin.conn.Close(); err != nil {
			lastErr = err
		}

		if err := plugin.process.Kill(); err != nil {
			lastErr = err
		}
	}

	pm.plugins = make(map[string]*managedPlugin)

	return lastErr
}

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
