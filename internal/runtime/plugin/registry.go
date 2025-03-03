// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package plugin

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"

	extism "github.com/extism/go-sdk"
	"github.com/tetratelabs/wazero"
)

type PluginRegistry struct {
	lookupDirs    []string
	loadedPlugins map[string]*LoadedPlugin
	mu            sync.RWMutex
}

type LoadedPlugin struct {
	wasmPath   string
	wasmPlugin *extism.Plugin
	manifest   *Manifest
}

func NewPluginRegistry(lookupDirs []string) *PluginRegistry {
	return &PluginRegistry{
		lookupDirs:    lookupDirs,
		loadedPlugins: make(map[string]*LoadedPlugin),
	}
}

func (pr *PluginRegistry) GetWASMPath(id string) (string, error) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	plugin, ok := pr.loadedPlugins[id]
	if !ok {
		return "", fmt.Errorf("failed to find plugin with id %s", id)
	}

	return plugin.wasmPath, nil
}

func (pr *PluginRegistry) GetWASMPlugin(id string) (*extism.Plugin, error) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	plugin, ok := pr.loadedPlugins[id]
	if !ok {
		return nil, fmt.Errorf("failed to find plugin with id %s", id)
	}

	return plugin.wasmPlugin, nil
}

func (pr *PluginRegistry) GetManifest(id string) (*Manifest, error) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	plugin, ok := pr.loadedPlugins[id]
	if !ok {
		return nil, fmt.Errorf("failed to find plugin with id %s", id)
	}

	return plugin.manifest, nil
}

func (pr *PluginRegistry) LoadPlugins() error {
	var wg sync.WaitGroup

	for _, lookupDir := range pr.lookupDirs {
		manifestPaths, err := FindManifestPaths(lookupDir)
		if err != nil {
			return fmt.Errorf("failed to load plugins: %s", err)
		}

		for _, manifestPath := range manifestPaths {
			wg.Add(1)

			go func() {
				defer wg.Done()

				if _, err := pr.LoadPlugin(manifestPath); err != nil {
					log.Printf("Error loading plugin: %v\n", err)
				}
			}()
		}
	}

	wg.Wait()

	return nil
}

func (pr *PluginRegistry) LoadPlugin(manifestPath string) (*LoadedPlugin, error) {
	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin: %w", err)
	}

	wasmPath := filepath.Join(filepath.Dir(manifestPath), strings.TrimPrefix(manifest.Entrypoint, "./"))

	wasmPlugin, err := loadWASMPlugin(manifest, wasmPath)
	if err != nil {
		return nil,
			fmt.Errorf("failed to load plugin: %w", err)
	}

	loadedPlugin := &LoadedPlugin{
		wasmPath:   wasmPath,
		wasmPlugin: wasmPlugin,
		manifest:   manifest,
	}

	pr.mu.Lock()
	pr.loadedPlugins[manifest.Id] = loadedPlugin
	pr.mu.Unlock()

	return loadedPlugin, nil
}

func loadWASMPlugin(m *Manifest, wasmPath string) (*extism.Plugin, error) {
	ctx := context.Background()
	cache := wazero.NewCompilationCache()
	defer cache.Close(ctx)

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Name: m.Id,
				Path: wasmPath,
			},
		},
	}

	config := extism.PluginConfig{
		EnableWasi:    true,
		ModuleConfig:  wazero.NewModuleConfig(),
		RuntimeConfig: wazero.NewRuntimeConfig().WithCompilationCache(cache),
	}

	plugin, err := extism.NewPlugin(ctx, manifest, config, []extism.HostFunction{})
	if err != nil {
		return nil, fmt.Errorf("failed to load WASM plugin: %v", err)
	}

	return plugin, nil
}
