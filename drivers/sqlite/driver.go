// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package sqlite

import (
	"context"
	"errors"

	"github.com/polyclient/polyclient/internal/database"
)

// Driver implements the database.Driver interface for SQLite.
type Driver struct{}

// Name implements database.Driver.Name().
func (d *Driver) Name() string {
	return "sqlite"
}

// Description implements database.Driver.Description().
func (d *Driver) Description() string {
	return "SQLite database driver for PolyClient"
}

// Version implements database.Driver.Version().
func (d *Driver) Version() string {
	return "1.0.0"
}

// Capabilities implements database.Driver.Capabilities().
func (d *Driver) Capabilities() database.DriverCapabilities {
	return database.DriverCapabilities{
		SupportsTransactions:        true,
		SupportsSchemaIntrospection: true,
		SupportsPreparedStatements:  true,
		SupportsStoredProcedures:    false,
		SupportsReplication:         false,
		SupportsTLS:                 false,
		SupportsSSH:                 false,
		IsDistributed:               false,
		IsEmbedded:                  true,
		MaxConnections:              1,
	}
}

// Connect implements database.Driver.Connect().
func (d *Driver) Connect(ctx context.Context, connectionString string) (database.Connection, error) {
	return nil, errors.New("not implemented")
}

// ValidateConnectionString implements database.Driver.ValidateConnectionString().
func (d *Driver) ValidateConnectionString(connectionString string) error {
	return errors.New("not implemented")
}

// Register registers the SQLite driver in the global database registry.
func Register() {
	registry := database.GetGlobalRegistry()
	_ = registry.RegisterDriver(&Driver{})
}
