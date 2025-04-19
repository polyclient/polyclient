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

	// Close terminates the connection.
	io.Closer
}
