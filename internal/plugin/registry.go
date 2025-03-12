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

// Registry manages discovery, loading, and execution of Wasm plugins.
type Registry struct {
	lookupDirs    []string
	loadedPlugins map[string]*LoadedPlugin
	mu            sync.RWMutex
}

// LoadedPlugin represents a loaded Wasm plugin.
type LoadedPlugin struct {
	wasmPath   string
	wasmPlugin *extism.Plugin
	manifest   *Manifest
}

// NewPluginRegistry initializes a registry with specified lookup directories.
func NewPluginRegistry(lookupDirs []string) *Registry {
	return &Registry{
		lookupDirs:    lookupDirs,
		loadedPlugins: map[string]*LoadedPlugin{},
	}
}

// GetWasmPath returns the Wasm file path for a plugin ID.
func (pr *Registry) GetWasmPath(id string) (string, error) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	plugin, ok := pr.loadedPlugins[id]
	if !ok {
		return "", fmt.Errorf("failed to find plugin with id %s", id)
	}

	return plugin.wasmPath, nil
}

// GetWasmPlugin returns the extism.Plugin instance for a plugin ID.
func (pr *Registry) GetWasmPlugin(id string) (*extism.Plugin, error) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	plugin, ok := pr.loadedPlugins[id]
	if !ok {
		return nil, fmt.Errorf("failed to find plugin with id %s", id)
	}

	return plugin.wasmPlugin, nil
}

// GetManifest returns the manifest for a plugin ID.
func (pr *Registry) GetManifest(id string) (*Manifest, error) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	plugin, ok := pr.loadedPlugins[id]
	if !ok {
		return nil, fmt.Errorf("failed to find plugin with id %s", id)
	}

	return plugin.manifest, nil
}

// CallFunction executes a function in the specified Wasm plugin.
func (pr *Registry) CallFunction(pluginID, functionName string, input []byte) ([]byte, error) {
	plugin, err := pr.GetWasmPlugin(pluginID)
	if err != nil {
		return nil, err
	}

	_, result, err := plugin.Call(functionName, input)
	if err != nil {
		return nil, fmt.Errorf("failed to call function %s: %w", functionName, err)
	}

	return result, nil
}

// CallFunctionWithContext executes a function in the specified Wasm plugin with a context.
func (pr *Registry) CallFunctionWithContext(ctx context.Context, pluginID, functionName string, input []byte) ([]byte, error) {
	plugin, err := pr.GetWasmPlugin(pluginID)
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
func (pr *Registry) LoadPlugins() error {
	var wg sync.WaitGroup

	for _, lookupDir := range pr.lookupDirs {
		manifestPaths, err := FindManifestPaths(lookupDir)
		if err != nil {
			return fmt.Errorf("failed to load plugins: %w", err)
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
func (pr *Registry) LoadPlugin(manifestPath string) (*LoadedPlugin, error) {
	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin: %w", err)
	}

	wasmPath := filepath.Join(filepath.Dir(manifestPath), strings.TrimPrefix(manifest.Entrypoint, "./"))

	wasmPlugin, err := loadWasmPlugin(manifest, wasmPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin: %w", err)
	}

	plugin := &LoadedPlugin{
		wasmPath,
		wasmPlugin,
		manifest,
	}

	pr.mu.Lock()
	pr.loadedPlugins[manifest.ID] = plugin
	pr.mu.Unlock()

	return plugin, nil
}

// UnloadPlugin unloads a single plugin from the registry.
func (pr *Registry) UnloadPlugin(id string) error {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	ctx := context.Background()

	plugin, ok := pr.loadedPlugins[id]
	if !ok {
		return fmt.Errorf("failed to find plugin with id %s", id)
	}

	if err := plugin.wasmPlugin.Close(ctx); err != nil {
		return fmt.Errorf("failed to close plugin: %w", err)
	}

	delete(pr.loadedPlugins, id)

	return nil
}

// loadWasmPlugin loads a single Wasm plugin from a file path.
func loadWasmPlugin(m *Manifest, wasmPath string) (*extism.Plugin, error) {
	ctx := context.Background()
	cache := wazero.NewCompilationCache()

	defer func() {
		if err := cache.Close(ctx); err != nil {
			log.Printf("Error closing cache: %v\n", err)
		}
	}()

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Name: m.ID,
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
		return nil, fmt.Errorf("failed to load Wasm plugin: %w", err)
	}

	return plugin, nil
}
