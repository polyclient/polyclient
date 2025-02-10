package plugin

import (
	"context"
	"fmt"
	"time"

	pb "github.com/polyclient/polyclient/proto"
)

// ExecuteAction executes the specified action on the plugin with the given name.
func (pm *PluginManager) Execute(pluginName, action string, payload []byte, metadata map[string]string) ([]byte, error) {
	pm.mu.Lock()
	plugin, exists := pm.plugins[pluginName]
	pm.mu.Unlock()

	if !exists {
		return nil, fmt.Errorf("plugin %s not found", pluginName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second) // TODO: make configurable
	defer cancel()

	result, err := plugin.client.Execute(ctx, &pb.PluginExecuteRequest{
		Action:  action,
		Payload: payload,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute action %s: %w", action, err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf("plugin failed to execute action %s", result.Error)
	}

	return result.Payload, nil
}
