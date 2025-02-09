package pluginsdk

// ActionHandler defines the action handlers for a plugin.
type ActionHandler func(payload []byte, metadata map[string]string) ([]byte, error)

// Plugin represents a generic plugin that can be used to extend any functionality of PolyClient.
type Plugin struct {
	Name     string
	Version  string
	Metadata map[string]string
	handlers map[string]ActionHandler
}

// NewPlugin creates a new plugin instance.
func NewPlugin(name, version string) *Plugin {
	return &Plugin{
		Name:     name,
		Version:  version,
		Metadata: make(map[string]string),
		handlers: make(map[string]ActionHandler),
	}
}

// AddMetadata adds a new metadata entry to the plugin.
func (p *Plugin) AddMetadata(key, value string) {
	p.Metadata[key] = value
}

// RegisterHandler registers a new action handler to the plugin.
func (p *Plugin) RegisterHandler(action string, handler ActionHandler) {
	p.handlers[action] = handler
}

// GetActions returns all registered actions of the plugin.
func (p *Plugin) GetActions() []string {
	actions := make([]string, 0, len(p.handlers))

	for action := range p.handlers {
		actions = append(actions, action)
	}

	return actions
}
