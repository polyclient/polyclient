// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package databasex

import "time"

// Schema represents a database schema.
type Schema struct {
	Name        string
	Tables      []Table
	Views       []View
	Procedures  []Procedure
	Functions   []Function
	Triggers    []Trigger
	Constraints []Constraint
	Sequences   []Sequence
	Types       []CustomType
}

// Table represents a database table.
type Table struct {
	Name          string
	Schema        string
	Columns       []Column
	Indexes       []Index
	PrimaryKey    *PrimaryKey
	ForeignKeys   []ForeignKey
	Constraints   []Constraint
	Comment       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	EstimatedRows int
	SizeBytes     int64
}

// Column represents a table's column.
type Column struct {
	Name            string
	Type            string
	Length          int
	Precision       int
	Scale           int
	DefaultValue    any
	IsNullable      bool
	IsPrimaryKey    bool
	IsUnique        bool
	IsIndexed       bool
	IsAutoIncrement bool
	IsForeignKey    bool
	ForeignKey      *ForeignKey
	Comment         string
	Position        int
	Collation       string
	CharacterSet    string
}

// PrimaryKey represents a table's primary key.
type PrimaryKey struct {
	Name        string
	Columns     []string
	IsClustered bool
	IndexType   string
}

// ForeignKey represents a table's foreign key.
type ForeignKey struct {
	Name       string
	Columns    []string
	RefSchema  string
	RefTable   string
	RefColumns []string
	OnDelete   string
	OnUpdate   string
}

// Index represents a database index.
type Index struct {
	Name        string
	Schema      string
	Table       string
	Columns     []IndexColumn
	IsUnique    bool
	IsClustered bool
	Type        string
	Filter      string
	Using       string
	SizeBytes   int
}

// IndexColumn represents a column in an index.
type IndexColumn struct {
	Name      string
	Position  int
	Direction string
	Length    int
}

// View represents a database view.
type View struct {
	Name           string
	Schema         string
	Columns        []Column
	Definition     string
	IsMaterialized bool
	Comment        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Procedure represents a database stored procedure.
type Procedure struct {
	Name       string
	Schema     string
	Parameters []Parameter
	ReturnType string
	Definition string
	Language   string
	Comment    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Function represents a database function.
type Function struct {
	Name       string
	Schema     string
	Parameters []Parameter
	ReturnType string
	Definition string
	Language   string
	Comment    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Parameter represents a stored procedure or function parameter.
type Parameter struct {
	Name         string
	Type         string
	Direction    string
	DefaultValue any
}

// Trigger represents a database trigger.
type Trigger struct {
	Name       string
	Schema     string
	Table      string
	Event      string
	Timing     string
	ForEachRow bool
	IsEnabled  bool
	Definition string
	Comment    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Constraint represents a database constraint.
type Constraint struct {
	Name       string
	Schema     string
	Table      string
	Type       string
	Definition string
	Columns    []string
	IsEnforced bool
	Comment    string
}

// Sequence represents a database sequence.
type Sequence struct {
	Name       string
	Schema     string
	DataType   string
	StartValue int
	Increment  int
	MinValue   int
	MaxValue   int
	CycleFlag  bool
	CacheSize  int
	Comment    string
}

// CustomType represents a user-defined data type.
type CustomType struct {
	Name       string
	Schema     string
	BaseType   string
	Definition string
	Comment    string
}

// Partitioning represents a table partitioning information.
type Partitioning struct {
	Type       string
	Expression string
	Columns    []string
	Partitions []Partition
}

// Partition represents a table partition.
type Partition struct {
	Name      string
	Value     string
	Bounds    string
	RowCount  int64
	SizeBytes int64
}

// SchemaChange represents a schema change operation.
type SchemaChange struct {
	Type         string
	Object       string
	ObjectType   string
	Schema       string
	SQL          string
	Description  string
	IsReversible bool
	IsRisky      bool
	Dependencies []SchemaChange
}
