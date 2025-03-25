// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package database

// DriverType represents the type of a database driver.
type DriverType string

const (
	// DriverTypeSQL represents SQL database driver type.
	DriverTypeSQL DriverType = "sql"

	// DriverTypeNoSQL represents NoSQL database driver type.
	DriverTypeNoSQL DriverType = "nosql"
)

// Driver represents a generic database driver.
type Driver[T Connection] interface {
	Name() string
	Type() DriverType
	CreateConnection(config ConnectionConfig) (T, error)
}

// DriverSQL represents an SQL database driver.
type DriverSQL interface {
	Driver[ConnectionSQL]
}

// DriverNoSQL represents a NoSQL database driver.
type DriverNoSQL interface {
	Driver[ConnectionNoSQL]
}

// AnyDriver is an alias for Driver with Connection type.
type AnyDriver = Driver[Connection]
