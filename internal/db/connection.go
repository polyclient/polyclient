// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package db

import (
	"context"
	"io"
)

// ConnectionConfig represents a driver-specific map of configuration options.
type ConnectionConfig map[string]any

// Connection represents a database connection.
type Connection interface {
	// Close terminates the connection.
	io.Closer

	// Ping verifies that the connection to the database is still alive.
	Ping(ctx context.Context) error
}

// ConnectionInfo contains metadata information about a database connection.
type ConnectionInfo interface {
	// ServerVersion returns the version of the database server used by the connection.
	ServerVersion(ctx context.Context) string

	// CurrentDatabase returns the name of the current database.
	CurrentDatabase(ctx context.Context) string
}
