// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package database

import "context"

// Pinger represents a database connection that can be pinged.
type Pinger interface {
	Ping() error
	PingContext(ctx context.Context) error
}
