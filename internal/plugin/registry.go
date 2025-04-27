// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package plugin

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bytecodealliance/wasmtime-go/v32"
)

// Registry manages discovery, loading, and execution of Wasm plugins.
type Registry struct {
	lookupDirs    []string
	loadedPlugins map[string]*LoadedPlugin
	engine        *wasmtime.Engine
	linker        *wasmtime.Linker
	mu            sync.RWMutex
}

// LoadedPlugin represents a loaded Wasm plugin.
type LoadedPlugin struct {
	wasmPath string
	instance *wasmtime.Instance
	store    *wasmtime.Store
	manifest *Manifest
}

// NewPluginRegistry initializes a registry with specified lookup directories.
func NewPluginRegistry(lookupDirs []string) (*Registry, error) {
	engine := wasmtime.NewEngine()
	linker := wasmtime.NewLinker(engine)

	err := linker.DefineWasi()
	if err != nil {
		return nil, fmt.Errorf("failed to define wasi: %w", err)
	}

	return &Registry{
		lookupDirs:    lookupDirs,
		loadedPlugins: map[string]*LoadedPlugin{},
		engine:        engine,
		linker:        linker,
	}, nil
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
func (pr *Registry) GetWasmPlugin(id string) (*wasmtime.Instance, *wasmtime.Store, error) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	plugin, ok := pr.loadedPlugins[id]
	if !ok {
		return nil, nil, fmt.Errorf("failed to find plugin with id %s", id)
	}

	return plugin.instance, plugin.store, nil
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
func (pr *Registry) CallFunction(pluginID, functionName, input string) ([]byte, error) {
	instance, store, err := pr.GetWasmPlugin(pluginID)
	if err != nil {
		return nil, err
	}

	funcExport := instance.GetFunc(store, functionName)
	if funcExport == nil {
		return nil, fmt.Errorf("function %s not found", functionName)
	}

	val, err := funcExport.Call(store, input)
	if err != nil {
		return nil, fmt.Errorf("failed to call function %s: %w", functionName, err)
	}

	result, ok := val.([]byte)
	if !ok {
		return nil, fmt.Errorf("unexpected result type for function %s", functionName)
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

	wasmPath := filepath.Join(
		filepath.Dir(manifestPath),
		strings.TrimPrefix(manifest.Entrypoint, "./"),
	)

	wasiConfig := wasmtime.NewWasiConfig()
	wasiConfig.InheritStdout()
	wasiConfig.InheritStderr()

	store := wasmtime.NewStore(pr.engine)
	store.SetWasi(wasiConfig)

	module, err := wasmtime.NewModuleFromFile(pr.engine, wasmPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin: %w", err)
	}

	wasmPlugin, err := pr.linker.Instantiate(store, module)
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin: %w", err)
	}

	plugin := &LoadedPlugin{
		wasmPath: wasmPath,
		instance: wasmPlugin,
		store:    store,
		manifest: manifest,
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

	_, ok := pr.loadedPlugins[id]
	if !ok {
		return fmt.Errorf("failed to find plugin with id %s", id)
	}

	delete(pr.loadedPlugins, id)

	return nil
}
