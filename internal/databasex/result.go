// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package databasex

import "time"

// Result represents the result of a database operation.
type Result interface {
	// LastInsertID returns the ID of the last inserted row.
	LastInsertID() (int64, error)

	// RowsAffected returns the number of rows affected by the operation.
	RowsAffected() (int64, error)

	// ExecutionTime returns how long the operation took to execute.
	ExecutionTime() time.Duration

	// Metadata returns the metadata associated with the result.
	Metadata() map[string]any
}

// Rows represents a set of database rows returned by a query.
type Rows interface {
	// Next prepares the next row for reading.
	Next() bool

	// Scan copies the columns in the current row into the values pointed at by dest.
	Scan(dest ...any) error

	// ColumnCount returns the number of columns in the result set.
	ColumnCount() int

	// ColumnNames returns the names of the columns in the result set.
	ColumnNames() ([]string, error)

	// ColumnTypes returns the types of the columns in the result set.
	ColumnTypes() ([]ColumnType, error)

	// Close closes the rows iterator.
	Close()

	// Err returns any error that occurred during iteration.
	Err() error
}

// ColumnType represents the type of a column in a database table.
type ColumnType struct {
	Name            string
	DatabaseType    string
	Length          int
	Precision       int
	Scale           int
	Nullable        bool
	IsPrimaryKey    bool
	IsAutoIncrement bool
	IsUnique        bool
	HasDefault      bool
	Default         any
}

// Paginator provides iterator-like behavior for a set of database rows.
type Paginator interface {
	// HasNext returns true if there are more pages.
	HasNext() bool

	// Next advances to the next page and returns its rows.
	Next() (Rows, error)

	// PageSize returns the number of rows per page.
	PageSize() int

	// CurrentPage returns the current page number (1-based).
	CurrentPage() int

	// TotalPages returns the total number of pages, if known.
	TotalPages() (int, bool)

	// TotalRows returns the total number of rows, if known.
	TotalRows() (int, bool)
}

type PaginationOptions struct {
	// PageSize is the number of rows per page.
	PageSize int

	// MaxPages is the maximum number of pages to fetch.
	MaxPages int
}
