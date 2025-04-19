// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package database

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// Registry manages database drivers.
type Registry[T Driver] struct {
	mu      sync.RWMutex
	drivers map[string]T
}

// NewRegistry creates a new database registry.
func NewRegistry[T Driver]() *Registry[T] {
	return &Registry[T]{
		drivers: map[string]T{},
	}
}

// Register registers a database driver.
func (r *Registry[T]) Register(driver T) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := strings.TrimSpace(driver.Name())
	if name == "" {
		return errors.New("driver name cannot be empty")
	}

	if _, exists := r.drivers[name]; exists {
		return fmt.Errorf("driver %s already registered", name)
	}

	r.drivers[name] = driver

	return nil
}

// Get returns a database driver by name.
func (r *Registry[T]) Get(name string) (T, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	driver, ok := r.drivers[name]

	return driver, ok
}

// List returns a list of all registered database drivers.
func (r *Registry[T]) List() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	drivers := make([]T, 0, len(r.drivers))
	for _, driver := range r.drivers {
		drivers = append(drivers, driver)
	}

	return drivers
}

// Unregister removes a driver if it exists.
func (r *Registry[T]) Unregister(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.drivers, name)
}

// Clear removes all drivers from the registry.
func (r *Registry[T]) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.drivers = map[string]T{}
}
