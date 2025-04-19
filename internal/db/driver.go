// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package db

import "context"

// Driver represents a generic database driver, capable of opening connections.
type Driver interface {
	// Name returns the unique name of the driver (e.g., "postgres", "mongodb").
	Name() string

	// Open connects to the database using the provided configuration.
	// Config should be a driver-specific map of configuration options.
	Open(ctx context.Context, config Config) (Connection, error)
}
