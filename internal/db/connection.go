// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package db

import (
	"context"
	"io"
)

// Connection represents an active database connection.
// The basic handle; specific capabilities are discovered via type assertions.
type Connection interface {
	// Ping verifies that the connection is alive.
	Ping(ctx context.Context) error

	// Info returns information about the connection.
	Info() ConnectionInfo

	// Close terminates the connection.
	io.Closer
}

// ConnectionInfo contains information about a database connection.
type ConnectionInfo interface {
	// ServerVersion returns the version of the database server used by the connection.
	ServerVersion(ctx context.Context) string

	// CurrentDatabase returns the name of the current database.
	CurrentDatabase(ctx context.Context) string
}
