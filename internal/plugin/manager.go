package plugin

import (
	"fmt"
	"os"
	"sync"

	"github.com/polyclient/polyclient/pkg/pluginsdk"
	pb "github.com/polyclient/polyclient/proto"
	"google.golang.org/grpc"
)

// PluginManager is a manager for PolyClient plugins.
type PluginManager struct {
	pluginsDirs []string
	plugins     map[string]*managedPlugin
	mu          sync.RWMutex
}

// managedPlugin represents a plugin managed by the PluginManager and its associated resources.
type managedPlugin struct {
	manifest pluginsdk.Manifest
	client   pb.PluginClient
	process  *os.Process
	conn     *grpc.ClientConn
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
