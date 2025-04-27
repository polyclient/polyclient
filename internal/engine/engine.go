// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/polyclient/polyclient/drivers/postgresql"
	"github.com/polyclient/polyclient/drivers/sqlite"
	"github.com/polyclient/polyclient/internal/db"
	"github.com/polyclient/polyclient/internal/env"
	"github.com/polyclient/polyclient/internal/plugin"
	"github.com/polyclient/polyclient/internal/settings"
	"github.com/polyclient/polyclient/internal/validator"
)

// Engine represents the application container for PolyClient.
// It orchestrates the various components of the application.
type Engine struct {
	Env             *env.Env                   // Environment manager for PolyClient
	Settings        *settings.Settings         // Application configuration
	DriversRegistry *db.Registry[db.Driver]    // Registry for database drivers
	PluginsRegistry *plugin.Registry           // Registry for plugins
	SDK             *db.SDK                    // SDK for database interactions
	Logger          *slog.Logger               // Logger for application logs
	Validator       *validator.CustomValidator // Validator for custom validation logic
}

// NewEngine creates a new Engine instance.
func NewEngine(ctx context.Context) (*Engine, error) {
	envMgr, err := env.NewManager()
	if err != nil {
		return nil, fmt.Errorf("Error setting up environment: %w", err)
	}

	logger, err := initLogger(envMgr)
	if err != nil {
		return nil, fmt.Errorf("Error initializing logger: %w", err)
	}

	logger.Info("Initializing PolyClient engine")

	stgs, err := initSettings(envMgr)
	if err != nil {
		return nil, fmt.Errorf("Error loading settings: %w", err)
	}

	driversReg, err := initDrivers(stgs)
	if err != nil {
		return nil, fmt.Errorf("Error loading database drivers: %w", err)
	}

	sdk, err := initSDK(envMgr, driversReg)
	if err != nil {
		return nil, fmt.Errorf("Error loading SDK: %w", err)
	}

	val, err := initValidator()
	if err != nil {
		return nil, fmt.Errorf("Error loading validator: %w", err)
	}

	pluginsReg, err := initPlugins(envMgr)
	if err != nil {
		logger.WarnContext(ctx, "Failed to load plugins. Continuing without plugins", "error", err)

		pluginsReg = nil
	}

	return &Engine{
		Env:             envMgr,
		Settings:        stgs,
		DriversRegistry: driversReg,
		PluginsRegistry: pluginsReg,
		SDK:             sdk,
		Logger:          logger,
		Validator:       val,
	}, nil
}

// initLogger configures the logger based on environment settings.
func initLogger(envMgr *env.Env) (*slog.Logger, error) {
	logLevelStr, err := envMgr.Get(env.PolyClientLogLevel)
	if err != nil {
		logLevelStr = "INFO"
	}

	var level slog.Level
	switch logLevelStr {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		return nil, fmt.Errorf("invalid log level: %s", logLevelStr)
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})

	return slog.New(handler), nil
}

// initSettings loads the application settings from environment variables.
func initSettings(envMgr *env.Env) (*settings.Settings, error) {
	settingsFile, err := envMgr.Get(env.PolyClientSettingsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to get PolyClient settings file: %w", err)
	}

	absPath, err := filepath.Abs(settingsFile)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %w", err)
	}

	settingsBytes, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PolyClient settings file: %w", err)
	}

	var s settings.Settings
	if err := json.Unmarshal(settingsBytes, &s); err != nil {
		return nil, fmt.Errorf("failed to unmarshal PolyClient settings: %w", err)
	}

	if err := s.Validate(); err != nil {
		return nil, err
	}

	return &s, nil
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
func initSDK(envMgr *env.Env, driversReg *db.Registry[db.Driver]) (*db.SDK, error) {
	connectionsDir, err := envMgr.Get(env.PolyClientConnectionsDir)
	if err != nil {
		return nil, fmt.Errorf("Error getting connections directory: %w", err)
	}

	connectionStore := db.NewFileConnectionStore(connectionsDir)
	connectionManager := db.NewConnectionManager(connectionStore, driversReg)

	return db.NewDatabaseSDK(connectionManager), nil
}

// initValidator loads a new CustomValidator instance.
func initValidator() (*validator.CustomValidator, error) {
	v, err := validator.NewCustomValidator()
	if err != nil {
		return nil, fmt.Errorf("Error creating validator: %w", err)
	}

	return v, nil
}

// initPlugins loads the built-in PolyClient plugins into the plugin registry.
// The built-in plugins are loaded from the POLYCLIENT_PLUGINS_DIR environment variable.
func initPlugins(envMgr *env.Env) (*plugin.Registry, error) {
	pluginsDir, err := envMgr.Get(env.PolyClientPluginsDir)
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
