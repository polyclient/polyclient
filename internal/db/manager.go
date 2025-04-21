package db

import (
	"context"
	"fmt"
	"log"
	"maps"
	"sync"
	"time"
)

// ConnectionManager manages active database connections.
type ConnectionManager struct {
	store         ConnectionStore
	registry      *Registry[Driver]
	activeConns   map[string]*ActiveConnection
	activeConnsMu sync.RWMutex
}

// ActiveConnection represents an active database connection.
type ActiveConnection struct {
	handle     Connection
	driverName string
}

// NewConnectionManager creates a new ConnectionManager instance.
func NewConnectionManager(store ConnectionStore, registry *Registry[Driver]) *ConnectionManager {
	return &ConnectionManager{
		store:       store,
		registry:    registry,
		activeConns: map[string]*ActiveConnection{},
	}
}

// GetStore returns the connection store.
func (m *ConnectionManager) GetStore() ConnectionStore {
	return m.store
}

// GetActiveConnection retrieves or establishes an active connection handle.
func (m *ConnectionManager) GetActiveConnection(ctx context.Context, name string) (*ActiveConnection, error) {
	m.activeConnsMu.Lock()
	defer m.activeConnsMu.Unlock()

	if active, ok := m.activeConns[name]; ok {
		pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		if err := active.handle.Ping(pingCtx); err == nil {
			return active, nil
		}

		log.Printf("WARN: Ping failed for existing active connection to %s. Attempting to reconnect.", name)

		go func(c Connection) { _ = c.Close() }(active.handle)
	}

	profile, err := m.store.GetProfile(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection profile: %w", err)
	}

	driver, ok := m.registry.Get(profile.Driver)
	if !ok {
		return nil, fmt.Errorf("driver %s not found", profile.Driver)
	}

	newConnHandle, err := driver.Open(ctx, profile.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := newConnHandle.Ping(pingCtx); err != nil {
		_ = newConnHandle.Close()
		return nil, fmt.Errorf("failed to ping connection: %w", err)
	}

	newActive := &ActiveConnection{
		handle:     newConnHandle,
		driverName: profile.Driver,
	}

	m.activeConns[name] = newActive

	log.Printf("INFO: Established and cached new active connection to %s for %s.", name, profile.Driver)

	return newActive, nil
}

// CloseActiveConnection explicitly closes and removes an active connection handle.
func (m *ConnectionManager) CloseActiveConnection(ctx context.Context, name string) bool {
	m.activeConnsMu.Lock()
	active, ok := m.activeConns[name]
	if ok {
		delete(m.activeConns, name)
	}
	m.activeConnsMu.Unlock()

	if ok {
		log.Printf("INFO: Explicitly closing active connection handle for %s", name)

		if err := active.handle.Close(); err != nil {
			log.Printf("Failed to close active connection handle for %s: %s", name, err)
		}

		return true
	}

	return false
}

// CloseAllActiveConnections closes all active connection handles during graceful shutdown.
func (m *ConnectionManager) CloseAllActiveConnections(ctx context.Context) error {
	m.activeConnsMu.Lock()
	connsToClose := make(map[string]*ActiveConnection, len(m.activeConns))
	maps.Copy(connsToClose, m.activeConns)
	m.activeConns = map[string]*ActiveConnection{}
	m.activeConnsMu.Unlock()

	var closeErrors []error
	var wg sync.WaitGroup

	for name, active := range connsToClose {
		wg.Add(1)

		go func() {
			defer wg.Done()
			log.Printf("INFO: Closing connection handle for %s", name)

			if err := active.handle.Close(); err != nil {
				m.activeConnsMu.Lock()
				closeErrors = append(closeErrors, fmt.Errorf("failed to close active connection handle for %s: %w", name, err))
				m.activeConnsMu.Unlock()
			}
		}()
	}

	wg.Wait()

	if len(closeErrors) > 0 {
		return fmt.Errorf("failed to close %d active connections: %v", len(closeErrors), closeErrors)
	}

	return nil
}

// StartHealthCheck periodically checks the liveness of active connections and removes stale ones.
func (m *ConnectionManager) StartHealthCheck(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.activeConnsMu.Lock()
			for name, active := range m.activeConns {
				if err := active.handle.Ping(ctx); err != nil {
					log.Printf("WARN: Health check failed for %s. Closing connection.", name)
					delete(m.activeConns, name)

					go func(c Connection) { _ = c.Close() }(active.handle)
				}
			}
			m.activeConnsMu.Unlock()
		}
	}
}
