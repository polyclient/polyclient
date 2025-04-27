// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package engine

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/polyclient/polyclient/drivers/postgresql"
	"github.com/polyclient/polyclient/drivers/sqlite"
	"github.com/polyclient/polyclient/internal/db"
	"github.com/polyclient/polyclient/internal/env"
	"github.com/polyclient/polyclient/internal/logger"
	"github.com/polyclient/polyclient/internal/plugin"
	"github.com/polyclient/polyclient/internal/settings"
	"github.com/polyclient/polyclient/internal/validator"
)

// Engine represents the application container for PolyClient.
// It orchestrates the various components of the application.
type Engine struct {
	Env             env.Provider               // Environment provider for PolyClient
	Settings        *settings.Settings         // Application configuration
	DriversRegistry *db.Registry[db.Driver]    // Registry for database drivers
	PluginsRegistry *plugin.Registry           // Registry for plugins
	SDK             *db.SDK                    // SDK for database interactions
	Logger          *logger.Logger             // Logger for application logs
	Validator       *validator.CustomValidator // Validator for custom validation logic
}

// NewEngine creates a new Engine instance.
func NewEngine(ctx context.Context) (*Engine, error) {
	provider := env.NewSystemProvider()
	if err := env.Prepare(provider); err != nil {
		return nil, fmt.Errorf("failed to initialize environment: %w", err)
	}

	stgs, err := initSettings(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to load settings: %w", err)
	}

	logger, err := logger.NewLogger(
		logger.WithFormat(parseLoggerFormat(stgs.Logging.Format)),
		logger.WithLevel(parseLoggerLevel(stgs.Logging.Level)),
		logger.WithRotationEnabled(stgs.Logging.Rotation.Enabled),
		logger.WithRotationMaxSize(stgs.Logging.Rotation.MaxSizeMB),
		logger.WithRotationMaxAge(stgs.Logging.Rotation.MaxAgeDays),
		logger.WithRotationMaxBackups(stgs.Logging.Rotation.MaxBackups),
		logger.WithRotationLocalTime(stgs.Logging.Rotation.LocalTime),
		logger.WithRotationCompress(stgs.Logging.Rotation.Compress),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}
	slog.SetDefault(logger.Logger)

	logger.Info("Initializing PolyClient engine")

	driversRegistry, err := initDrivers(stgs)
	if err != nil {
		return nil, fmt.Errorf("failed to load database drivers: %w", err)
	}

	sdk, err := initSDK(provider, driversRegistry)
	if err != nil {
		return nil, fmt.Errorf("failed to load database SDK: %w", err)
	}

	val, err := initValidator()
	if err != nil {
		return nil, fmt.Errorf("failed to load validator: %w", err)
	}

	// pluginsRegistry, err := initPlugins(provider)
	// if err != nil {
	// 	logger.WarnContext(ctx, "Failed to load plugins. Continuing without plugins", "error", err)

	// 	pluginsRegistry = nil
	// }

	return &Engine{
		Env:             provider,
		Settings:        stgs,
		DriversRegistry: driversRegistry,
		PluginsRegistry: nil,
		SDK:             sdk,
		Logger:          logger,
		Validator:       val,
	}, nil
}

// initSettings loads the application settings.
func initSettings(envProvider env.Provider) (*settings.Settings, error) {
	return settings.LoadFromEnv(envProvider)
}

// initDrivers initializes and registers the built-in database drivers (e.g., SQLite, PostgreSQL)
// into the provided driver registry. It returns an error if any driver registration fails.
func initDrivers(stgs *settings.Settings) (*db.Registry[db.Driver], error) {
	registry := db.NewRegistry[db.Driver]()

	if stgs.Drivers.SQLite.Enabled {
		if err := registry.Register(sqlite.NewDriver()); err != nil {
			return nil, fmt.Errorf("failed to register SQLite driver: %w", err)
		}
	}

	if stgs.Drivers.PostgreSQL.Enabled {
		if err := registry.Register(postgresql.NewDriver()); err != nil {
			return nil, fmt.Errorf("failed to register PostgreSQL driver: %w", err)
		}
	}

	return registry, nil
}

// initSDK loads the PolyClient SDK from the connections directory.
func initSDK(envProvider env.Provider, driversReg *db.Registry[db.Driver]) (*db.SDK, error) {
	connectionsDir, err := envProvider.Get(env.VariableConnectionsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to getting connections directory: %w", err)
	}

	connectionStore := db.NewFileConnectionStore(connectionsDir)
	connectionManager := db.NewConnectionManager(connectionStore, driversReg)

	return db.NewDatabaseSDK(connectionManager), nil
}

// initValidator loads a new CustomValidator instance.
func initValidator() (*validator.CustomValidator, error) {
	v, err := validator.NewCustomValidator()
	if err != nil {
		return nil, fmt.Errorf("failed to creating validator: %w", err)
	}

	return v, nil
}

// initPlugins loads the built-in PolyClient plugins into the plugin registry.
// The built-in plugins are loaded from the POLYCLIENT_PLUGINS_DIR environment variable.
func initPlugins(envProvider env.Provider) (*plugin.Registry, error) {
	pluginsDir, err := envProvider.Get(env.VariablePluginsDir)
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

// parseLoggerFormat parses the log format from a string.
func parseLoggerFormat(format string) logger.Format {
	switch format {
	case "json":
		return logger.JSONFormat
	case "text":
		return logger.TextFormat
	default:
		return logger.JSONFormat
	}
}

// parseLoggerLevel parses the log level from a string.
func parseLoggerLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
