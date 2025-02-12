package pluginsdk

import (
	"errors"
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
	Config   *Config
	Handlers map[string]ActionHandler
}

// NewPlugin creates a new plugin instance for the current executable.
//
// It assumes that the plugin config file is in the same directory as the
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
			return nil, errors.New("could not determine plugin directory")
		}

		pluginDir = filepath.Dir(currentFile)
	} else {
		// We're running a compiled binary
		pluginDir = filepath.Dir(execPath)
	}

	var configFilePath string

	for _, configFile := range supportedConfigFiles {
		if _, err := os.Stat(filepath.Join(pluginDir, configFile)); err == nil {
			configFilePath = filepath.Join(pluginDir, configFile)
			break
		}
	}

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to find plugin config in %s", pluginDir)
	}

	config, err := LoadConfig(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin config from %s: %w", configFilePath, err)
	}

	return &Plugin{
		Config:   config,
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
