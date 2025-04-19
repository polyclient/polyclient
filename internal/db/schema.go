package db

import "time"

// ObjectType represents the type of a database object (Table, View, etc.)
type ObjectType string

const (
	ObjectTypeDatabase         ObjectType = "DATABASE"
	ObjectTypeSchema           ObjectType = "SCHEMA"
	ObjectTypeTable            ObjectType = "TABLE"
	ObjectTypeView             ObjectType = "VIEW"
	ObjectTypeMaterializedView ObjectType = "MATERIALIZED_VIEW"
	ObjectTypeIndex            ObjectType = "INDEX"
	ObjectTypeSequence         ObjectType = "SEQUENCE"
	ObjectTypeFunction         ObjectType = "FUNCTION"
	ObjectTypeProcedure        ObjectType = "PROCEDURE"
	ObjectTypeTrigger          ObjectType = "TRIGGER"
	ObjectTypeType             ObjectType = "TYPE"
	ObjectTypeConstraint       ObjectType = "CONSTRAINT"
)

// ConstraintType identifies the type of constraint.
type ConstraintType string

const (
	ConstraintTypePrimaryKey ConstraintType = "PK"
	ConstraintTypeForeignKey ConstraintType = "FK"
	ConstraintTypeUnique     ConstraintType = "UK"
	ConstraintTypeCheck      ConstraintType = "CHECK"
)

// GenericObjectIdentifier is a generic identifier for a database object.
type GenericObjectIdentifier struct {
	Schema string `json:"schema,omitempty"` // Schema name (can be empty if object type doesn't live in a schema or for DB-level objects)
	Name   string `json:"name"`             // Object name
}

// GenericObjectOwner is a generic owner for a database object.
type GenericObjectOwner struct {
	Owner string `json:"owner,omitempty"` // Object owner/user if applicable
}

// GenericObjectComment is a generic comment for a database object.
type GenericObjectComment struct {
	Comment string `json:"comment,omitempty"` // Object comment/description
}

// GenericObjectDefinition is a generic definition for a database object.
type GenericObjectDefinition struct {
	Definition string `json:"definition,omitempty"` // Object definition (e.g., CREATE TABLE statement)
}

// GenericObjectTimestamps defines timestamps for a database object.
type GenericObjectTimestamps struct {
	CreatedAt time.Time `json:"createdAt,omitempty"` // Timestamp when the object was created
}

// GenericObjectAttributes is a generic set of attributes for a database object.
type GenericObjectAttributes struct {
	Attributes map[string]any `json:"attributes,omitempty"` // Driver-specific attributes (e.g., engine, charset, collation, etc.)
}

// --- Summary Structs (for Lists requests) ---

// DatabaseSummary provides basic information about a database.
type DatabaseSummary struct {
	GenericObjectIdentifier
	GenericObjectOwner
	GenericObjectAttributes
}

// SchemaSummary provides basic information about a schema.
type SchemaSummary struct {
	GenericObjectIdentifier
	GenericObjectOwner
	GenericObjectAttributes

	Catalog string `json:"catalog,omitempty"` // The database name if different from schema context
}

// TableSummary provides basic information about a table.
type TableSummary struct {
	GenericObjectIdentifier
	GenericObjectOwner
	GenericObjectAttributes

	IsTemporary bool `json:"isTemporary,omitempty"` // True if the table is a temporary table
}

// ViewSummary provides basic information about a view.
type ViewSummary struct {
	GenericObjectIdentifier
	GenericObjectOwner
	GenericObjectAttributes
}

// MaterializedViewSummary provides basic information about a materialized view.
type MaterializedViewSummary struct {
	ViewSummary
}

// ColumnSummary provides basic information about a column.
type ColumnSummary struct {
	GenericObjectIdentifier
	GenericObjectOwner
	GenericObjectAttributes
}

// IndexSummary provides basic information about an index.
type IndexSummary struct {
	GenericObjectIdentifier
	GenericObjectOwner
	GenericObjectAttributes

	IsUnique  bool `json:"isUnique,omitempty"` // True if the index is unique
	IsPrimary bool `json:"isPrimary,omitempty"`
}

// SequenceSummary provides basic information about a sequence.
type SequenceSummary struct {
	GenericObjectIdentifier
	GenericObjectOwner
	GenericObjectAttributes
}

// FunctionSummary provides basic information about a function.
type FunctionSummary struct {
	GenericObjectIdentifier
	GenericObjectOwner
	GenericObjectAttributes
}

// ProcedureSummary provides basic information about a procedure.
type ProcedureSummary struct {
	GenericObjectIdentifier
	GenericObjectOwner
	GenericObjectAttributes
}

// ArgumentSummary provides basic information about an argument.
type ArgumentSummary struct {
	GenericObjectIdentifier
	GenericObjectOwner
	GenericObjectAttributes
}

// TriggerSummary provides basic information about a trigger.
type TriggerSummary struct {
	GenericObjectIdentifier
	GenericObjectOwner
	GenericObjectAttributes
}

// TypeSummary provides basic information about a type.
type TypeSummary struct {
	GenericObjectIdentifier
	GenericObjectOwner
	GenericObjectAttributes

	Type     string `json:"type"`               // TYPE or DOMAIN
	BaseType string `json:"baseType,omitempty"` // Base SQL type if applicable (e.g., VARCHAR)
}

// ConstraintSummary provides basic information about a constraint.
type ConstraintSummary struct {
	GenericObjectIdentifier
	GenericObjectOwner
	GenericObjectAttributes

	Type         ConstraintType `json:"type"` // Constraint type (e.g., "PK", "FK", "UK", "CHECK")
	TargetSchema string         `json:"targetSchema"`
	TargetTable  string         `json:"targetTable"`
	IsEnabled    bool           `json:"isEnabled"` // True if the constraint is enforced
}

// --- Detail Structs (for Get requests) ---

// DatabaseDetail provides detailed information about a database.
type DatabaseDetail struct {
	DatabaseSummary
	GenericObjectComment
	GenericObjectTimestamps
}

// SchemaDetail provides detailed information about a schema.
type SchemaDetail struct {
	SchemaSummary
	GenericObjectComment
	GenericObjectTimestamps
}

// TableDetail provides detailed information about a table.
type TableDetail struct {
	TableSummary
	GenericObjectComment
	GenericObjectTimestamps

	Columns     []ColumnDetail     `json:"columns,omitempty"`     // List of table columns
	Indexes     []IndexDetail      `json:"indexes,omitempty"`     // List of table indexes
	Constraints []ConstraintDetail `json:"constraints,omitempty"` // List of table constraints
}

// ViewDetail provides detailed information about a view.
type ViewDetail struct {
	ViewSummary
	GenericObjectComment
	GenericObjectTimestamps
	GenericObjectDefinition

	Columns []ColumnDetail `json:"columns,omitempty"` // List of view columns
}

// MaterializedViewDetail represents a database materialized view.
type MaterializedViewDetail struct {
	MaterializedViewSummary
	GenericObjectComment
	GenericObjectTimestamps
	GenericObjectDefinition

	Columns []ColumnDetail `json:"columns,omitempty"` // List of materialized view columns
}

// ColumnDetail provides detailed information about a column.
type ColumnDetail struct {
	ColumnSummary
	GenericObjectComment

	Position     int    `json:"position"`               // Ordinal position (1-based)
	BaseType     string `json:"baseType,omitempty"`     // Simplified base type (e.g., "VARCHAR", "NUMERIC") - driver best effort
	FullDataType string `json:"fullDataType"`           // Database-specific full type (e.g., "VARCHAR(100)", "NUMERIC(10,2)")
	IsNullable   bool   `json:"isNullable"`             // True if the column allows NULL values
	IsPrimaryKey bool   `json:"isPrimaryKey,omitempty"` // True if the column is a primary key
	DefaultValue any    `json:"defaultValue,omitempty"` // Default value for the column (e.g., 7, "now()", etc.)
}

// IndexDetail provides detailed information about an index.
type IndexDetail struct {
	IndexSummary
	GenericObjectComment
	GenericObjectTimestamps
	GenericObjectDefinition

	Columns        []IndexColumnDetail `json:"columns"`                  // Columns included in the index key
	IncludeColumns []string            `json:"includeColumns,omitempty"` // Columns included via INCLUDE clause
	Definition     string              `json:"definition,omitempty"`     // CREATE INDEX statement / expression for expression index
	Predicate      *string             `json:"predicate,omitempty"`      // Partial index condition (WHERE clause)
}

// IndexColumnDetail describes a column within an index key.
// Note: Does NOT embed ObjectIdentifier.
type IndexColumnDetail struct {
	Name       string         `json:"name"`                 // Column name or expression
	Position   int            `json:"position"`             // Position in the index key (1-based)
	SortOrder  string         `json:"sortOrder,omitempty"`  // "ASC", "DESC"
	NullsOrder string         `json:"nullsOrder,omitempty"` // "NULLS FIRST", "NULLS LAST"
	Attrs      map[string]any `json:"attrs,omitempty"`      // e.g., Collation, OperatorClass
}

// ConstraintDetail provides detailed information about a constraint.
type ConstraintDetail struct {
	ConstraintSummary
	GenericObjectComment
	GenericObjectDefinition

	// FK specific fields (only populated if Type == ConstraintTypeForeignKey)
	ForeignTableSchema string   `json:"foreignTableSchema,omitempty"` // Foreign table schema (e.g., "dbo")
	ForeignTableName   string   `json:"foreignTableName,omitempty"`   // Foreign table name
	ForeignColumns     []string `json:"foreignColumns,omitempty"`     // Foreign table column names
	OnUpdateAction     string   `json:"onUpdateAction,omitempty"`     // NO ACTION, RESTRICT, CASCADE, SET NULL, SET DEFAULT
	OnDeleteAction     string   `json:"onDeleteAction,omitempty"`     // NO ACTION, RESTRICT, CASCADE, SET NULL
	MatchOption        string   `json:"matchOption,omitempty"`        // SIMPLE, FULL, PARTIAL
}

// SequenceDetail provides detailed information about a sequence.
type SequenceDetail struct {
	SequenceSummary
	GenericObjectComment
	GenericObjectTimestamps

	DataType     string               `json:"dataType"`          // Database-specific data type (e.g., "SERIAL", "BIGSERIAL", etc.)
	StartValue   int64                `json:"startValue"`        // Starting value of the sequence
	IncrementBy  int64                `json:"incrementBy"`       // Increment value of the sequence
	MinValue     int64                `json:"minValue"`          // Minimum value of the sequence
	MaxValue     int64                `json:"maxValue"`          // Maximum value of the sequence
	CurrentValue int64                `json:"currentValue"`      // Current value of the sequence
	IsCyclic     bool                 `json:"isCyclic"`          // True if the sequence is cyclic
	CacheSize    int64                `json:"cacheSize"`         // Cache size of the sequence
	OwnedBy      *SequenceOwnerDetail `json:"ownedBy,omitempty"` // If linked to a SERIAL/IDENTITY column
}

type SequenceOwnerDetail struct {
	SchemaName string `json:"schemaName"`
	TableName  string `json:"tableName"`
	ColumnName string `json:"columnName"`
}

// TypeDetail provides comprehensive UDT/Domain information.
type TypeDetail struct {
	TypeSummary
	GenericObjectComment
	GenericObjectTimestamps
	GenericObjectDefinition

	EnumLabels   []string `json:"enumLabels,omitempty"`   // For ENUM types
	IsNullable   bool     `json:"isNullable"`             // True if the type allows NULL values
	DefaultValue any      `json:"defaultValue,omitempty"` // Default value for the type (e.g., 7, "now()", etc.)
}

// FunctionDetail provides detailed information about a function.
type FunctionDetail struct {
	FunctionSummary
	GenericObjectDefinition

	Language  string           `json:"language"`  // Function language (e.g., "SQL", "PL/SQL", etc.)
	Arguments []ArgumentDetail `json:"arguments"` // List of input/output arguments
}

// ProcedureDetail provides detailed information about a procedure.
type ProcedureDetail struct {
	ProcedureSummary
	GenericObjectDefinition

	Language  string           `json:"language,omitempty"`  // Procedure language (e.g., "SQL", "PL/SQL", etc.)
	Arguments []ArgumentDetail `json:"arguments,omitempty"` // List of input/output arguments
}

// ArgumentDetail provides detailed information about an argument.
type ArgumentDetail struct {
	ArgumentSummary

	Position int    `json:"position"`       // Ordinal position (1-based)
	Mode     string `json:"mode,omitempty"` // Argument mode (e.g., "IN", "OUT", "INOUT")
	DataType string `json:"dataType"`       // Database-specific data type (e.g., "VARCHAR", "INTEGER", "TIMESTAMP")
}

// TriggerDetail provides detailed information about a trigger.
type TriggerDetail struct {
	TriggerSummary
	GenericObjectDefinition

	TargetTableName  string   `json:"targetTable"`  // Name of the table that the trigger is associated with
	TargetSchemaName string   `json:"targetSchema"` // Name of the schema of the table that the trigger is associated with
	Events           []string `json:"events"`       // List of events that trigger the trigger (e.g., "INSERT", "UPDATE", "DELETE")
	Timing           string   `json:"timing"`       // Trigger timing (e.g., "BEFORE", "AFTER", "INSTEAD OF")
	Orientation      string   `json:"orientation"`  // Trigger orientation (e.g., "FOR EACH ROW", "FOR EACH STATEMENT")
}
