package pluginsdk

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// ActionHandler defines the action handlers for a plugin.
type ActionHandler func(payload []byte) ([]byte, error)

// Plugin represents a generic plugin that can be used to extend any functionality of PolyClient.
type Plugin struct {
	Manifest *Manifest
	Handlers map[string]ActionHandler
}

// NewPlugin creates a new plugin instance for the current executable.
//
// It assumes that the manifest.json file is in the same directory as the
// executable. If the file is not found, it returns an error.
func NewPlugin() (*Plugin, error) {
	var pluginDir string

	execPath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("could not determine executable path: %w", err)
	}

	// Check if we're running with 'go run' (executable will be in /tmp)
	if strings.HasPrefix(execPath, os.TempDir()) {
		// We're in development mode
		_, currentFile, _, ok := runtime.Caller(1)
		if !ok {
			return nil, fmt.Errorf("could not determine plugin directory")
		}
		pluginDir = filepath.Dir(currentFile)
	} else {
		// We're running a compiled binary
		pluginDir = filepath.Dir(execPath)
	}

	manifestPath := filepath.Join(pluginDir, "manifest.json")

	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("manifest.json not found in plugin directory: %s", pluginDir)
	}

	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin manifest at %s: %w", manifestPath, err)
	}

	return &Plugin{
		Manifest: manifest,
		Handlers: make(map[string]ActionHandler),
	}, nil
}

// RegisterHandler registers a new action handler to the plugin.
func (p *Plugin) RegisterHandler(action string, handler ActionHandler) {
	p.Handlers[action] = handler
}

// GetActions returns all registered actions of the plugin.
func (p *Plugin) GetActions() []string {
	actions := make([]string, 0, len(p.Handlers))

	for action := range p.Handlers {
		actions = append(actions, action)
	}

	return actions
}
