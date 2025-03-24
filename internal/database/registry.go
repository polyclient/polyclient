// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package database

import (
	"fmt"
	"sync"
)

// Registry manages database drivers.
type Registry struct {
	mu      sync.RWMutex
	drivers map[string]Driver
	sql     map[string]SQLDriver
	// nosql   map[string]NoSQLDriver
}

// newRegistry creates a new database registry.
func newRegistry() *Registry {
	return &Registry{
		drivers: make(map[string]Driver),
		sql:     make(map[string]SQLDriver),
		// nosql:   make(map[string]NoSQLDriver),
	}
}

// RegisterDriver registers a database driver.
func (r *Registry) RegisterDriver(driver Driver) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := driver.Name()
	if _, exists := r.drivers[name]; exists {
		return fmt.Errorf("driver %s already registered", name)
	}

	r.drivers[name] = driver

	// Also register in the appropriate specific map
	if sqlDriver, ok := driver.(SQLDriver); ok {
		r.sql[name] = sqlDriver
	}
	// if noSQLDriver, ok := driver.(NoSQLDriver); ok {
	// 	r.nosql[name] = noSQLDriver
	// }

	return nil
}

// GetDriver returns a database driver by name.
func (r *Registry) GetDriver(name string) (Driver, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	driver, ok := r.drivers[name]

	return driver, ok
}

// GetSQLDriver returns a SQL database driver by name.
func (r *Registry) GetSQLDriver(name string) (SQLDriver, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	driver, ok := r.sql[name]

	return driver, ok
}

// GetNoSQLDriver returns a NoSQL database driver by name.
// func (r *Registry) GetNoSQLDriver(name string) (NoSQLDriver, bool) {
// 	r.mu.RLock()
// 	defer r.mu.RUnlock()

// 	driver, ok := r.nosql[name]
// 	return driver, ok
// }

// ListDrivers returns a list of all registered database drivers.
func (r *Registry) ListDrivers() []Driver {
	r.mu.RLock()
	defer r.mu.RUnlock()

	drivers := make([]Driver, 0, len(r.drivers))
	for _, driver := range r.drivers {
		drivers = append(drivers, driver)
	}

	return drivers
}

// ListSQLDrivers returns a list of all registered SQL database drivers.
func (r *Registry) ListSQLDrivers() []SQLDriver {
	r.mu.Lock()
	defer r.mu.Unlock()

	drivers := make([]SQLDriver, 0, len(r.sql))
	for _, driver := range r.sql {
		drivers = append(drivers, driver)
	}

	return drivers
}

// ListNoSQLDrivers returns a list of all registered NoSQL database drivers.
// func (r *Registry) ListNoSQLDrivers() []NoSQLDriver {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()

// 	drivers := make([]NoSQLDriver, 0, len(r.nosql))
// 	for _, driver := range r.nosql {
// 		drivers = append(drivers, driver)
// 	}
// 	return drivers
// }

// UnregisterDriver unregisters a database driver by name.
func (r *Registry) UnregisterDriver(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.drivers, name)
	delete(r.sql, name)
	// delete(r.nosql, name)
}

// Clear removes all drivers from the registry.
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.drivers = make(map[string]Driver)
	r.sql = make(map[string]SQLDriver)
	// r.nosql = make(map[string]NoSQLDriver)
}

// global registry instance.
var (
	globalRegistry *Registry
	once           sync.Once
)

// GetGlobalRegistry returns the global database registry.
func GetGlobalRegistry() *Registry {
	once.Do(func() {
		globalRegistry = newRegistry()
	})

	return globalRegistry
}
