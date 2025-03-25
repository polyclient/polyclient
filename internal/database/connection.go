// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package database

import "context"

// Connection represents a generic database connection.
type Connection interface {
	Close() error
	Ping() error
	PingContext(ctx context.Context) error
}

type ConnectionConfig interface {
	Validate() error
}

// ConnectionSQL represents an SQL database connection.
type ConnectionSQL interface {
	Connection
}

// ConnectionNoSQL represents a NoSQL database connection.
type ConnectionNoSQL interface {
	Connection
}
