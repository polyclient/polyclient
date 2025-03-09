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

// PluginRegistry manages discovery, loading, and execution of WASM plugins.
type PluginRegistry struct {
	lookupDirs    []string
	loadedPlugins map[string]*LoadedPlugin
	mu            sync.RWMutex
}

// LoadedPlugin represents a loaded WASM plugin.
type LoadedPlugin struct {
	wasmPath   string
	wasmPlugin *extism.Plugin
	manifest   *PluginManifest
}

// NewPluginRegistry initializes a registry with specified lookup directories.
func NewPluginRegistry(lookupDirs []string) *PluginRegistry {
	return &PluginRegistry{
		lookupDirs:    lookupDirs,
		loadedPlugins: make(map[string]*LoadedPlugin),
	}
}

// GetWASMPath returns the WASM file path for a plugin ID.
func (pr *PluginRegistry) GetWASMPath(id string) (string, error) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	plugin, ok := pr.loadedPlugins[id]
	if !ok {
		return "", fmt.Errorf("failed to find plugin with id %s", id)
	}

	return plugin.wasmPath, nil
}

// GetWASMPlugin returns the extism.Plugin instance for a plugin ID.
func (pr *PluginRegistry) GetWASMPlugin(id string) (*extism.Plugin, error) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	plugin, ok := pr.loadedPlugins[id]
	if !ok {
		return nil, fmt.Errorf("failed to find plugin with id %s", id)
	}

	return plugin.wasmPlugin, nil
}

// GetManifest returns the manifest for a plugin ID.
func (pr *PluginRegistry) GetManifest(id string) (*PluginManifest, error) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	plugin, ok := pr.loadedPlugins[id]
	if !ok {
		return nil, fmt.Errorf("failed to find plugin with id %s", id)
	}

	return plugin.manifest, nil
}

// CallFunction executes a function in the specified WASM plugin.
func (pr *PluginRegistry) CallFunction(pluginId, functionName string, input []byte) ([]byte, error) {
	plugin, err := pr.GetWASMPlugin(pluginId)
	if err != nil {
		return nil, err
	}

	_, result, err := plugin.Call(functionName, input)
	if err != nil {
		return nil, fmt.Errorf("failed to call function %s: %w", functionName, err)
	}

	return result, nil
}

// CallFunctionWithContext executes a function in the specified WASM plugin with a context.
func (pr *PluginRegistry) CallFunctionWithContext(ctx context.Context, pluginId, functionName string, input []byte) ([]byte, error) {
	plugin, err := pr.GetWASMPlugin(pluginId)
	if err != nil {
		return nil, err
	}

	_, result, err := plugin.CallWithContext(ctx, functionName, input)
	if err != nil {
		return nil, fmt.Errorf("failed to call function %s: %w", functionName, err)
	}

	return result, nil
}

// LoadPlugins scans lookup directories and loads available plugins.
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

// LoadPlugin loads a single plugin from a manifest file path.
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

	plugin := &LoadedPlugin{
		wasmPath,
		wasmPlugin,
		manifest,
	}

	pr.mu.Lock()
	pr.loadedPlugins[manifest.Id] = plugin
	pr.mu.Unlock()

	return plugin, nil
}

// UnloadPlugin unloads a single plugin from the registry.
func (pr *PluginRegistry) UnloadPlugin(id string) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	ctx := context.Background()

	plugin, ok := pr.loadedPlugins[id]
	if !ok {
		return fmt.Errorf("failed to find plugin with id %s", id)
	}

	plugin.wasmPlugin.Close(ctx)

	delete(pr.loadedPlugins, id)

	return nil
}

// loadWASMPlugin loads a single WASM plugin from a file path.
func loadWASMPlugin(m *PluginManifest, wasmPath string) (*extism.Plugin, error) {
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
