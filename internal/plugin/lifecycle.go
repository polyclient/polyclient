package plugin

// Shutdown stops all plugin processes and closes all plugin connections.
func (pm *PluginManager) Shutdown() error {
	return pm.cleanup()
}
