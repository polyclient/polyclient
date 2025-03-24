// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package database

import (
	"context"
	"time"
)

// IsolationLevel represents the isolation level for a transaction.
type IsolationLevel int

const (
	// IsolationLevelDefault uses the database's default isolation level.
	IsolationLevelDefault IsolationLevel = iota

	// IsolationLevelReadUncommitted allows reading uncommitted changes.
	IsolationLevelReadUncommitted

	// IsolationLevelReadCommitted allows reading committed changes.
	IsolationLevelReadCommitted

	// IsolationLevelRepeatableRead ensures the same data will be read if re-queried.
	IsolationLevelRepeatableRead

	// IsolationLevelSerializable provides the highest isolation level.
	IsolationLevelSerializable
)

// Transaction represents a database transaction.
type Transaction interface {
	// Query executes a query within a transaction.
	Query(ctx context.Context, query string, args ...any) (Rows, error)

	// Execute executes a statement within a transaction.
	Execute(ctx context.Context, query string, args ...any) (Result, error)

	// Commit commits the transaction.
	Commit() error

	// Rollback rolls back the transaction.
	Rollback() error

	// Savepoint creates a savepoint within the transaction.
	Savepoint(name string) error

	// RollbackToSavepoint rolls back to a savepoint within the transaction.
	RollbackToSavepoint(name string) error

	// IsActive returns whether the transaction is still active.
	IsActive() bool
}

// TransactionOptions holds the options for a transaction.
type TransactionOptions struct {
	// IsReadOnly specifies whether the transaction should be read-only.
	IsReadOnly bool

	// IsolationLevel specifies the isolation level for the transaction.
	IsolationLevel IsolationLevel

	// Name provides an optional name for the transaction.
	Name string

	// Timeout specifies the maximum duration for the transaction.
	Timeout time.Duration
}

// IsolationLevelName returns the string name of the isolation level.
func IsolationLevelName(level IsolationLevel) string {
	switch level {
	case IsolationLevelReadUncommitted:
		return "READ UNCOMMITTED"
	case IsolationLevelReadCommitted:
		return "READ COMMITTED"
	case IsolationLevelRepeatableRead:
		return "REPEATABLE READ"
	case IsolationLevelSerializable:
		return "SERIALIZABLE"
	default:
		return "DEFAULT"
	}
}
