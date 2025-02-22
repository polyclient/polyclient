// Copyright (C) 2025 Juan Mesa and contributors
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3
// as published by the Free Software Foundation, with the Runtime
// Library Exception. See the COPYING.RUNTIME file for details.

package plugin

import (
	"context"
	"fmt"
	"log"

	extism "github.com/extism/go-sdk"
	"github.com/tetratelabs/wazero"
)

type PluginManager struct {
	loadedPlugins []*extism.Plugin
}

func NewPluginManager() *PluginManager {
	return &PluginManager{}
}

func (m *PluginManager) LoadPlugins() error {
	ctx := context.Background()

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{Path: "./plugins/sqlite/polyclient-sqlite.wasm"},
		},
	}

	config := extism.PluginConfig{
		EnableWasi:    true,
		ModuleConfig:  wazero.NewModuleConfig(),
		RuntimeConfig: wazero.NewRuntimeConfig(),
	}

	plugin, err := extism.NewPlugin(ctx, manifest, config, []extism.HostFunction{})
	if err != nil {
		return fmt.Errorf("failed to load plugin: %v", err)
	}

	result, output, err := plugin.Call("greet", []byte("Juan"))
	if err != nil {
		return fmt.Errorf("failed to call greet: %v", err)
	}

	log.Println(result, string(output))

	m.loadedPlugins = append(m.loadedPlugins, plugin)

	return nil
}
