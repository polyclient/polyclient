// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package db

import "io"

// Result represents a generic database result set.
type Result interface {
	// Next advances to the next row/document/item in the result set.
	// Returns false when iteration is complete or an error occurs.
	Next() bool

	// Scan copies the columns in the current row/document/item into
	// the values pointed at by dest.
	Scan(dest ...any) error

	// Err returns any error that occurred during iteration (e.g., by Next).
	Err() error

	// Close releases any resources associated with the result set.
	io.Closer

	// Columns returns the column names, if applicable (primarily for tabular data).
	Columns() ([]string, error)

	// RowsAffected returns the number of rows affected by an INSERT/UPDATE/DELETE, if available.
	RowsAffected() (int64, error) // Returns 0/-1 and/or error if not applicable/supported
}
