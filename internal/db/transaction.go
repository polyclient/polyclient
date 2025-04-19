package db

import "context"

// Transactional defines the capability to manage transactions.
type Transactional interface {
	// BeginTx starts a new transaction with the given options.
	BeginTx(ctx context.Context, opts *TxOptions) (Tx, error)
}

// Tx represents an active database transaction.
type Tx interface {
	// Commit commits the transaction.
	Commit() error

	// Rollback aborts the transaction.
	Rollback() error

	QueryExecutor
}

// TxOptions represents options for starting a transaction.
type TxOptions struct {
	IsolationLevel TxIsolationLevel
	AccessMode     TxAccessMode
}

// TxIsolationLevel represents a transaction isolation level.
type TxIsolationLevel int

const (
	IsolationLevelDefault TxIsolationLevel = iota
	IsolationLevelReadUncommitted
	IsolationLevelReadCommitted
	IsolationLevelRepeatableRead
	IsolationLevelSnapshot
	IsolationLevelSerializable
)

// TxAccessMode represents a transaction access mode.
type TxAccessMode int

const (
	AccessModeDefault TxAccessMode = iota
	AccessModeReadOnly
	AccessModeReadWrite
)
