// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package database

// Connection represents a generic database connection.
type Connection interface {
	Close() error
}

// ConnectionSQL represents a generic SQL database connection.
type ConnectionSQL interface {
	Connection
	Pinger
	SQLQuerier
}

// ConnectionNoSQL represents a generic NoSQL database connection.
type ConnectionNoSQL interface {
	Connection
	Pinger
}
