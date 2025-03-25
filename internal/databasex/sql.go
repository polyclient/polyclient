// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package databasex

import (
	"context"
)

// SQLSchemaInspector handles schema inspection operations.
type SQLSchemaInspector interface {
	// GetSchema returns database schema information.
	GetSchema(ctx context.Context, conn Connection) (Schema, error)
}

// SQLQueryExplainer handles query execution plan operations.
type SQLQueryExplainer interface {
	// ExplainQuery returns the query execution plan.
	ExplainQuery(ctx context.Context, conn Connection, query string) (string, error)
}

// SQLStatementPreparer handles prepared statements operations.
type SQLStatementPreparer interface {
	// PrepareStatement creates a prepared statement.
	PrepareStatement(ctx context.Context, conn Connection, query string) (SQLPreparedStatement, error)
}

// SQLDialectProvider handles SQL dialect operations.
type SQLDialectProvider interface {
	// GetDialect returns information about the SQL dialect.
	GetDialect() SQLDialect

	// QuoteIdentifier quotes an identifier according to the SQL dialect.
	QuoteIdentifier(identifier string) string

	// GetDataTypes returns the data types supported by the SQL dialect.
	GetDataTypes() []SQLDataType
}

// SQLDriver extends Driver with SQL-specific functionality.
type SQLDriver interface {
	Driver
	SQLSchemaInspector
	SQLQueryExplainer
	SQLStatementPreparer
	SQLDialectProvider
}

// SQLConnection extends Connection with SQL-specific functionality.
type SQLConnection interface {
	Connection

	// GetTables returns a list of tables belonging to the specified schema.
	GetTables(ctx context.Context, schema string) ([]Table, error)

	// GetColumns returns a list of columns belonging to the specified table.
	GetColumns(ctx context.Context, schema, table string) ([]Column, error)

	// GetViews returns a list of views belonging to the specified schema.
	GetViews(ctx context.Context, schema string) ([]View, error)
}

// SQLPreparedStatement represents a prepared SQL statement.
type SQLPreparedStatement interface {
	// Execute executes the prepared statement.
	Execute(ctx context.Context, args ...any) (Result, error)

	// Query executes the prepared statement and returns rows.
	Query(ctx context.Context, args ...any) (Rows, error)

	// Close closes the prepared statement.
	Close() error

	// NumInput returns the number of placeholders parameters.
	NumInput() int
}

// SQLDialect represents an SQL dialect.
type SQLDialect struct {
	// Name of the dialect (e.g., "mysql", "postgres", "sqlite").
	Name string

	// IdentifierQuote is the character used to quote identifiers.
	IdentifierQuote string

	// StringLiteralQuote is the character used to quote string literals.
	StringLiteralQuote string

	// SupportsLimitOffset indicates whether the dialect supports LIMIT/OFFSET.
	SupportsLimitOffset bool

	// SupportsReturning indicates whether the dialect supports RETURNING.
	SupportsReturning bool

	// SupportsCTE indicates whether the dialect supports Common Table Expressions (CTEs).
	SupportsCTE bool

	// SupportWindowFunctions indicates whether the dialect supports window functions.
	SupportsWindowFunctions bool

	// SupportsUpsert indicates whether the dialect supports UPSERT operations.
	SupportsUpsert bool
}

// SQLExecutionPlan represents a query execution plan.
type SQLExecutionPlan struct {
	PlanMode     string
	Operation    string
	Relation     string
	Alias        string
	StartupCost  float64
	TotalCost    float64
	Rows         int
	Width        int
	FilteredRows float64
	ActualTime   float64
	ActualRows   int
	Loops        int
	SubPlans     []SQLExecutionPlan
	Details      map[string]any
}

// SQLDataType represents an SQL data type.
type SQLDataType struct {
	// Name is the name of the data type.
	Name string

	// Aliases are alternative names for the data type.
	Aliases []string

	// Category is the broad category of the data type (numeric, string, etc.).
	Category string

	// HasLength indicates whether the data type has a length.
	HasLength bool

	// HasPrecision indicates whether the data type has a precision.
	HasPrecision bool

	// IsNumeric indicates whether the data type is numeric.
	IsNumeric bool

	// DefaultLength is the default length for the data type.
	DefaultLength int

	// DefaultPrecision is the default precision for the data type.
	DefaultPrecision int

	// DefaultScale is the default scale for the data type.
	DefaultScale int
}
