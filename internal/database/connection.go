// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package database

// Connection represents a generic database connection.
type Connection interface {
	Close() error
}

type ConnectionSQL interface {
	Connection
	Pinger
	SQLQuerier
}

type ConnectionNoSQL interface {
	Connection
	Pinger
}
