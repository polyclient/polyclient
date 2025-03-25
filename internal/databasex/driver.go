// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package databasex

import "context"

// DriverType represents the type of a database driver.
type DriverType string

const (
	// DriverTypeSQL represents SQL database drivers.
	DriverTypeSQL DriverType = "sql"

	// DriverTypeNoSQL represents NoSQL database drivers.
	DriverTypeNoSQL DriverType = "nosql"

	// DriverTypeGraph represents graph database drivers.
	DriverTypeGraph DriverType = "graph"

	// DriverTypeTimeSeries represents time series database drivers.
	DriverTypeTimeSeries DriverType = "timeseries"

	// DriverTypeKeyValue represents key-value database drivers.
	DriverTypeKeyValue DriverType = "keyvalue"
)

// DriverDescriptor provides metadata about a driver.
type DriverDescriptor interface {
	// Name returns the name of the driver (e.g., "mysql", "postgres", "sqlite").
	Name() string

	// Description returns a human-readable description of the driver.
	Description() string

	// Version returns the version of the driver.
	Version() string

	// Type returns the type of the driver.
	Type() DriverType
}

// DriverCapabilityProvider provides driver capabilities.
type DriverCapabilityProvider interface {
	// Capabilities returns the capabilities of the driver.
	Capabilities() DriverCapabilities
}

// DriverConnectionValidator provides methods to validate connections.
type DriverConnectionValidator interface {
	// ValidateConnectionString checks if the connection string is valid.
	ValidateConnectionString(connectionString string) error
}

// DriverConnector handles connection operations.
type DriverConnector interface {
	// Connect returns a connection to the database.
	Connect(ctx context.Context, connectionString string) (Connection, error)
}

// Driver represents a generic database driver interface.
type Driver interface {
	DriverDescriptor
	DriverCapabilityProvider
	DriverConnectionValidator
	DriverConnector
}

// DriverCapabilities represents the capabilities supported by a database driver.
type DriverCapabilities struct {
	SupportsTransactions        bool
	SupportsSchemaIntrospection bool
	SupportsPreparedStatements  bool
	SupportsStoredProcedures    bool
	SupportsReplication         bool
	SupportsTLS                 bool
	SupportsSSH                 bool
	IsDistributed               bool
	IsEmbedded                  bool
	MaxConnections              int
}
