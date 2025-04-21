package application

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/polyclient/polyclient/drivers/postgresql"
	"github.com/polyclient/polyclient/drivers/sqlite"
	"github.com/polyclient/polyclient/internal/config"
	"github.com/polyclient/polyclient/internal/db"
	"github.com/polyclient/polyclient/internal/env"
	"github.com/polyclient/polyclient/internal/plugin"
	"github.com/polyclient/polyclient/internal/validator"
)

// App is the application container.
type Application struct {
	Config          *config.PolyClientConfig
	DriversRegistry *db.Registry[db.Driver]
	PluginsRegistry *plugin.Registry
	SDK             *db.SDK
	Logger          *slog.Logger
	Validator       *validator.CustomValidator
}

// NewApplication creates a new application container.
func NewApplication(ctx context.Context) (*Application, error) {
	if err := env.GetManager().Setup(); err != nil {
		return nil, fmt.Errorf("Error setting up environment: %w", err)
	}

	driversRegistry, err := loadDrivers()
	if err != nil {
		return nil, fmt.Errorf("Error loading database drivers: %w", err)
	}

	connectionsDir, err := env.GetManager().Get(env.EnvPolyClientConnectionsDir)
	if err != nil {
		return nil, fmt.Errorf("Error getting connections directory: %w", err)
	}

	connectionStore := db.NewFileConnectionStore(connectionsDir)
	connectionManager := db.NewConnectionManager(connectionStore, driversRegistry)
	sdk := db.NewDatabaseSDK(connectionManager)

	// pluginsRegistry, err := loadPlugins()
	// if err != nil {
	// 	return nil, fmt.Errorf("Error loading plugins:", err)
	// }

	logger := slog.Default()

	v, err := validator.NewCustomValidator()
	if err != nil {
		return nil, fmt.Errorf("Error creating validator: %w", err)
	}

	return &Application{
		Config:          nil,
		DriversRegistry: driversRegistry,
		PluginsRegistry: nil,
		SDK:             sdk,
		Logger:          logger,
		Validator:       v,
	}, nil
}

// loadDrivers loads the built-in database drivers into the driver registry.
func loadDrivers() (*db.Registry[db.Driver], error) {
	registry := db.NewRegistry[db.Driver]()

	if err := registry.Register(sqlite.NewDriver()); err != nil {
		return nil, fmt.Errorf("failed to register SQLite driver: %w", err)
	}

	if err := registry.Register(postgresql.NewDriver()); err != nil {
		return nil, fmt.Errorf("failed to register PostgreSQL driver: %w", err)
	}

	return registry, nil
}

// loadPlugins loads the built-in PolyClient plugins into the plugin registry.
// The built-in plugins are loaded from the POLYCLIENT_PLUGINS_DIR environment variable.
func loadPlugins() (*plugin.Registry, error) {
	pluginsDir, err := env.GetManager().Get(env.EnvPolyClientPluginsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get PolyClient plugins directory: %w", err)
	}

	lookupPaths := []string{pluginsDir}

	pr, err := plugin.NewPluginRegistry(lookupPaths)
	if err != nil {
		return nil, fmt.Errorf("failed to create plugin registry: %w", err)
	}

	if err := pr.LoadPlugins(); err != nil {
		return nil, fmt.Errorf("failed to load plugins: %w", err)
	}

	return pr, nil
}
