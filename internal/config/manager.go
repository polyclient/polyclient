// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package config

// Manager represents a PolyClient configuration manager.
type Manager struct {
	defaults       *PolyClientConfig
	loadedConfig   *PolyClientConfig
	schemaPath     string
	userConfigPath string
}

// NewManager creates a new PolyClient configuration manager.
func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) SetDefaults(defaults *PolyClientConfig) {
	m.defaults = defaults
}

func (m *Manager) SetSchemaPath(schemaPath string) {
	m.schemaPath = schemaPath
}

func (m *Manager) SetUserConfigPath(userConfigPath string) {
	m.userConfigPath = userConfigPath
}

func (m *Manager) SetLoadedConfig(loadedConfig *PolyClientConfig) {
	m.loadedConfig = loadedConfig
}

func (m *Manager) GetDefaults() *PolyClientConfig {
	return m.defaults
}

func (m *Manager) GetSchemaPath() string {
	return m.schemaPath
}

func (m *Manager) GetUserConfigPath() string {
	return m.userConfigPath
}

func (m *Manager) GetLoadedConfig() *PolyClientConfig {
	return m.loadedConfig
}

func (m *Manager) GetConfig() *PolyClientConfig {
	return m.loadedConfig
}

func (m *Manager) SetConfig(config *PolyClientConfig) {
	m.loadedConfig = config
}

func (m *Manager) SaveConfig() error {
	return nil
}
