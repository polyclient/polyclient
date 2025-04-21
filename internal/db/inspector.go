// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package db

import "context"

// DatabaseLister defines the capability to list databases in the database.
type DatabaseLister interface {
	// ListDatabases returns a list of all databases in the database.
	ListDatabases(ctx context.Context) ([]DatabaseSummary, error)
}

// DatabaseGetter defines the capability to retrieve a specific database by name.
type DatabaseGetter interface {
	// GetDatabase returns the database with the given name.
	GetDatabase(ctx context.Context, name string) (*DatabaseDetail, error)
}

// SchemaLister defines the capability to list schemas in the database.
type SchemaLister interface {
	// ListSchemas returns a list of all schemas in the database.
	ListSchemas(ctx context.Context) ([]SchemaSummary, error)
}

// SchemaGetter defines the capability to retrieve a specific schema by name.
type SchemaGetter interface {
	// GetSchema returns the schema with the given name.
	GetSchema(ctx context.Context, name string) (*SchemaDetail, error)
}

// TableLister defines the capability to list tables in the database.
type TableLister interface {
	// ListTables returns a list of all tables in the database.
	ListTables(ctx context.Context, options ...ListTablesOption) ([]TableSummary, error)
}

// TableGetter defines the capability to retrieve a specific table by name.
type TableGetter interface {
	// GetTable returns the table with the given name.
	GetTable(ctx context.Context, name string, options ...GetTableOption) (*TableDetail, error)
}

// ViewLister defines the capability to list views in the database.
type ViewLister interface {
	// ListViews returns a list of all views in the database.
	ListViews(ctx context.Context) ([]ViewSummary, error)
}

// ViewGetter defines the capability to retrieve a specific view by name.
type ViewGetter interface {
	// GetView returns the view with the given name.
	GetView(ctx context.Context, name string) (*ViewDetail, error)
}

// MaterializedViewLister defines the capability to list materialized views in the database.
type MaterializedViewLister interface {
	// ListMaterializedViews returns a list of all materialized views in the database.
	ListMaterializedViews(ctx context.Context) ([]MaterializedViewSummary, error)
}

// MaterializedViewGetter defines the capability to retrieve a specific materialized view by name.
type MaterializedViewGetter interface {
	// GetMaterializedView returns the materialized view with the given name.
	GetMaterializedView(ctx context.Context, name string) (*MaterializedViewDetail, error)
}

// ColumnLister defines the capability to list columns in the database.
type ColumnLister interface {
	// ListColumns returns a list of all columns in the database.
	ListColumns(ctx context.Context) ([]ColumnSummary, error)
}

// ColumnGetter defines the capability to retrieve a specific column by name.
type ColumnGetter interface {
	// GetColumn returns the column with the given name.
	GetColumn(ctx context.Context, name string) (*ColumnDetail, error)
}

// IndexLister defines the capability to list indexes in the database.
type IndexLister interface {
	// ListIndexes returns a list of all indexes in the database.
	ListIndexes(ctx context.Context) ([]IndexSummary, error)
}

// IndexGetter defines the capability to retrieve a specific index by name.
type IndexGetter interface {
	// GetIndex returns the index with the given name.
	GetIndex(ctx context.Context, name string) (*IndexDetail, error)
}

// ConstraintLister defines the capability to list constraints in the database.
type ConstraintLister interface {
	// ListConstraints returns a list of all constraints in the database.
	ListConstraints(ctx context.Context) ([]ConstraintSummary, error)
}

// ConstraintGetter defines the capability to retrieve a specific constraint by name.
type ConstraintGetter interface {
	// GetConstraint returns the constraint with the given name.
	GetConstraint(ctx context.Context, name string) (*ConstraintDetail, error)
}

// SequenceLister defines the capability to list sequences in the database.
type SequenceLister interface {
	// ListSequences returns a list of all sequences in the database.
	ListSequences(ctx context.Context) ([]SequenceSummary, error)
}

// SequenceGetter defines the capability to retrieve a specific sequence by name.
type SequenceGetter interface {
	// GetSequence returns the sequence with the given name.
	GetSequence(ctx context.Context, name string) (*SequenceDetail, error)
}

// TypeLister defines the capability to list user-defined types in the database.
type TypeLister interface {
	// ListTypes returns a list of all types in the database.
	ListTypes(ctx context.Context) ([]TypeSummary, error)
}

// TypeGetter defines the capability to retrieve a specific user-defined type by name.
type TypeGetter interface {
	// GetType returns the type with the given name.
	GetType(ctx context.Context, name string) (*TypeDetail, error)
}

// FunctionLister defines the capability to list functions in the database.
type FunctionLister interface {
	// ListFunctions returns a list of all functions in the database.
	ListFunctions(ctx context.Context) ([]FunctionSummary, error)
}

// FunctionGetter defines the capability to retrieve a specific function by name.
type FunctionGetter interface {
	// GetFunction returns the function with the given name.
	GetFunction(ctx context.Context, name string) (*FunctionDetail, error)
}

// ProcedureLister defines the capability to list procedures in the database.
type ProcedureLister interface {
	// ListProcedures returns a list of all procedures in the database.
	ListProcedures(ctx context.Context) ([]ProcedureSummary, error)
}

// ProcedureGetter defines the capability to retrieve a specific procedure by name.
type ProcedureGetter interface {
	// GetProcedure returns the procedure with the given name.
	GetProcedure(ctx context.Context, name string) (*ProcedureDetail, error)
}

// ArgumentLister defines the capability to list arguments (e.g., of functions/procedures) in the database.
type ArgumentLister interface {
	// ListArguments returns a list of all arguments in the database.
	ListArguments(ctx context.Context) ([]ArgumentSummary, error)
}

// ArgumentGetter defines the capability to retrieve a specific argument by name.
type ArgumentGetter interface {
	// GetArgument returns the argument with the given name.
	GetArgument(ctx context.Context, name string) (*ArgumentDetail, error)
}

// TriggerLister defines the capability to list triggers in the database.
type TriggerLister interface {
	// ListTriggers returns a list of all triggers in the database.
	ListTriggers(ctx context.Context) ([]TriggerSummary, error)
}

// TriggerGetter defines the capability to retrieve a specific trigger by name.
type TriggerGetter interface {
	// GetTrigger returns the trigger with the given name.
	GetTrigger(ctx context.Context, name string) (*TriggerDetail, error)
}
