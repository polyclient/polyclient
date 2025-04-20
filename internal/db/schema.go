package db

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
	ObjectTypeUnknown          ObjectType = "UNKNOWN"
)

// ConstraintType identifies the type of constraint.
type ConstraintType string

const (
	ConstraintTypePrimaryKey ConstraintType = "PRIMARY KEY"
	ConstraintTypeForeignKey ConstraintType = "FOREIGN KEY"
	ConstraintTypeUnique     ConstraintType = "UNIQUE"
	ConstraintTypeCheck      ConstraintType = "CHECK"
)

// ConstraintActionType identifies the type of constraint action.
type ConstraintActionType string

const (
	ConstraintActionTypeNoAction   ConstraintActionType = "NO ACTION"
	ConstraintActionTypeRestrict   ConstraintActionType = "RESTRICT"
	ConstraintActionTypeCascade    ConstraintActionType = "CASCADE"
	ConstraintActionTypeSetNull    ConstraintActionType = "SET NULL"
	ConstraintActionTypeSetDefault ConstraintActionType = "SET DEFAULT"
)

// ConstraintMatchType identifies the type of constraint match.
type ConstraintMatchType string

const (
	ConstraintMatchTypeFull    ConstraintMatchType = "FULL"
	ConstraintMatchTypePartial ConstraintMatchType = "PARTIAL"
	ConstraintMatchTypeSimple  ConstraintMatchType = "SIMPLE"
)

// TriggerTiming identifies the type of trigger timing.
type TriggerTiming string

const (
	TriggerTimingBefore    TriggerTiming = "BEFORE"
	TriggerTimingAfter     TriggerTiming = "AFTER"
	TriggerTimingInsteadOf TriggerTiming = "INSTEAD OF"
)

// TriggerEvent identifies the type of trigger event.
type TriggerEvent string

const (
	TriggerEventInsert   TriggerEvent = "INSERT"
	TriggerEventUpdate   TriggerEvent = "UPDATE"
	TriggerEventDelete   TriggerEvent = "DELETE"
	TriggerEventTruncate TriggerEvent = "TRUNCATE"
)

// TriggerOrientation identifies the type of trigger orientation.
type TriggerOrientation string

const (
	TriggerOrientationRow       TriggerOrientation = "ROW"
	TriggerOrientationStatement TriggerOrientation = "STATEMENT"
)

// TriggerEnabled identifies the type of trigger enabled.
type TriggerEnabled string

const (
	TriggerEnabledOrigin   TriggerEnabled = "ORIGIN"
	TriggerEnabledAlways   TriggerEnabled = "ALWAYS"
	TriggerEnabledReplica  TriggerEnabled = "REPLICA"
	TriggerEnabledDisabled TriggerEnabled = "DISABLED"
)

// GenericObjectIdentifier is a generic identifier for a database object.
type GenericObjectIdentifier struct {
	Schema string `json:"schema,omitempty"` // Schema name (can be empty if object type doesn't live in a schema or for DB-level objects)
	Name   string `json:"name"`             // Object name
	OID    int    `json:"oid,omitempty"`    // Object ID (if applicable)
}

// GenericObjectComment is a generic comment for a database object.
type GenericObjectComment struct {
	Comment string `json:"comment,omitempty"` // Object comment/description
}

// GenericObjectDefinition is a generic definition for a database object.
type GenericObjectDefinition struct {
	Definition string `json:"definition,omitempty"` // Object definition (e.g., CREATE TABLE statement)
}

// GenericObjectAttributes is a generic set of attributes for a database object.
type GenericObjectAttributes struct {
	Attributes map[string]any `json:"attributes,omitempty"` // Driver-specific attributes (e.g., engine, charset, collation, etc.)
}

// --- Summary Structs (for Lists requests) ---

// DatabaseSummary provides basic information about a database.
type DatabaseSummary struct {
	GenericObjectIdentifier
	GenericObjectAttributes
}

// SchemaSummary provides basic information about a schema.
type SchemaSummary struct {
	GenericObjectIdentifier
	GenericObjectAttributes

	Catalog string `json:"catalog,omitempty"` // The database name if different from schema context
}

// TableSummary provides basic information about a table.
type TableSummary struct {
	GenericObjectIdentifier
	GenericObjectAttributes

	IsTemporary bool `json:"isTemporary,omitempty"` // True if the table is a temporary table
}

// ViewSummary provides basic information about a view.
type ViewSummary struct {
	GenericObjectIdentifier
	GenericObjectAttributes
}

// MaterializedViewSummary provides basic information about a materialized view.
type MaterializedViewSummary struct {
	ViewSummary
}

// ColumnSummary provides basic information about a column.
type ColumnSummary struct {
	GenericObjectIdentifier
	GenericObjectAttributes
}

// IndexSummary provides basic information about an index.
type IndexSummary struct {
	GenericObjectIdentifier
	GenericObjectAttributes

	IsUnique  bool `json:"isUnique,omitempty"` // True if the index is unique
	IsPrimary bool `json:"isPrimary,omitempty"`
}

// SequenceSummary provides basic information about a sequence.
type SequenceSummary struct {
	GenericObjectIdentifier
	GenericObjectAttributes
}

// FunctionSummary provides basic information about a function.
type FunctionSummary struct {
	GenericObjectIdentifier
	GenericObjectAttributes
}

// ProcedureSummary provides basic information about a procedure.
type ProcedureSummary struct {
	GenericObjectIdentifier
	GenericObjectAttributes
}

// ArgumentSummary provides basic information about an argument.
type ArgumentSummary struct {
	GenericObjectIdentifier
	GenericObjectAttributes
}

// TriggerSummary provides basic information about a trigger.
type TriggerSummary struct {
	GenericObjectIdentifier
	GenericObjectAttributes
}

// TypeSummary provides basic information about a type.
type TypeSummary struct {
	GenericObjectIdentifier
	GenericObjectAttributes

	Type     string `json:"type"`               // TYPE or DOMAIN
	BaseType string `json:"baseType,omitempty"` // Base SQL type if applicable (e.g., VARCHAR)
}

// ConstraintSummary provides basic information about a constraint.
type ConstraintSummary struct {
	GenericObjectIdentifier
	GenericObjectAttributes

	Type ConstraintType `json:"type"` // Constraint type (e.g., "PRIMARY KEY", "FOREIGN KEY", "UNIQUE", "CHECK", etc.)
}

// --- Detail Structs (for Get requests) ---

// DatabaseDetail provides detailed information about a database.
type DatabaseDetail struct {
	DatabaseSummary
	GenericObjectComment
}

// SchemaDetail provides detailed information about a schema.
type SchemaDetail struct {
	SchemaSummary
	GenericObjectComment
}

// TableDetail provides detailed information about a table.
type TableDetail struct {
	TableSummary
	GenericObjectComment

	Columns     []ColumnDetail     `json:"columns,omitempty"`     // List of table columns
	Indexes     []IndexDetail      `json:"indexes,omitempty"`     // List of table indexes
	Constraints []ConstraintDetail `json:"constraints,omitempty"` // List of table constraints
	Triggers    []TriggerDetail    `json:"triggers,omitempty"`    // List of table triggers
}

// ViewDetail provides detailed information about a view.
type ViewDetail struct {
	ViewSummary
	GenericObjectComment
	GenericObjectDefinition

	Columns []ColumnDetail `json:"columns,omitempty"` // List of view columns
}

// MaterializedViewDetail represents a database materialized view.
type MaterializedViewDetail struct {
	MaterializedViewSummary
	GenericObjectComment

	GenericObjectDefinition

	Columns []ColumnDetail `json:"columns,omitempty"` // List of materialized view columns
}

// ColumnDetail provides detailed information about a column.
type ColumnDetail struct {
	ColumnSummary
	GenericObjectComment

	TypeName         string `json:"typeName"`                 // DB-specific type name (e.g., int4, varchar, timestamptz, etc.)
	TypeFamily       string `json:"typeFamily"`               // Logical type category (e.g., integer, character varying, timestamp with time zone, etc.)
	OrdinalPosition  int    `json:"ordinalPosition"`          // Ordinal position of the column
	MaxLength        int    `json:"maxLength,omitempty"`      // Max length (for char/varchar/text fields)
	Precision        int    `json:"precision,omitempty"`      // Numeric precision (for decimals)
	Scale            int    `json:"scale,omitempty"`          // Numeric scale (for decimals)
	IsNullable       bool   `json:"isNullable"`               // True if the column allows NULL values
	IsSelfReferenced bool   `json:"isSelfReferenced"`         // True if the column is a self-referenced foreign key
	IsGenerated      bool   `json:"isGenerated"`              // True if the column is generated by the database
	IsUpdatable      bool   `json:"isUpdatable"`              // True if the column is updatable  (Columns in base tables are always updatable, columns in views not necessarily)
	GenerationExpr   string `json:"generationExpr,omitempty"` // SQL expression for computed/generated columns
	DefaultValue     any    `json:"defaultValue,omitempty"`   // Default value for the column (e.g., 7, "now()", etc.)
}

// IndexDetail provides detailed information about an index.
type IndexDetail struct {
	IndexSummary
	GenericObjectComment
	GenericObjectDefinition

	Columns    string `json:"columns"`              // Columns included in the index key
	Definition string `json:"definition,omitempty"` // CREATE INDEX statement / expression for expression index
}

// ConstraintDetail provides detailed information about a constraint.
type ConstraintDetail struct {
	ConstraintSummary
	GenericObjectComment
	GenericObjectDefinition

	// FK specific fields (only populated if Type == ConstraintTypeForeignKey)
	ForeignSchema  string               `json:"foreignSchema,omitempty"`  // Foreign table schema (e.g., "public")
	ForeignTable   string               `json:"foreignTable,omitempty"`   // Foreign table name (e.g., "planet")
	OnUpdateAction ConstraintActionType `json:"onUpdateAction,omitempty"` // NO ACTION, RESTRICT, CASCADE, SET NULL, SET DEFAULT
	OnDeleteAction ConstraintActionType `json:"onDeleteAction,omitempty"` // NO ACTION, RESTRICT, CASCADE, SET NULL
	MatchOption    ConstraintMatchType  `json:"matchOption,omitempty"`    // SIMPLE, FULL, PARTIAL
}

// SequenceDetail provides detailed information about a sequence.
type SequenceDetail struct {
	SequenceSummary
	GenericObjectComment

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

// SequenceOwnerDetail provides information about the owner of a sequence.
type SequenceOwnerDetail struct {
	SchemaName string `json:"schemaName"`
	TableName  string `json:"tableName"`
	ColumnName string `json:"columnName"`
}

// TypeDetail provides comprehensive UDT/Domain information.
type TypeDetail struct {
	TypeSummary
	GenericObjectComment
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

	OrdinalPosition int    `json:"ordinalPosition"` // Ordinal position in the argument list
	Mode            string `json:"mode,omitempty"`  // Argument mode (e.g., "IN", "OUT", "INOUT")
	DataType        string `json:"dataType"`        // Database-specific data type (e.g., "VARCHAR", "INTEGER", "TIMESTAMP")
}

// TriggerDetail provides detailed information about a trigger.
type TriggerDetail struct {
	TriggerSummary
	GenericObjectComment
	GenericObjectDefinition

	Timing         TriggerTiming      `json:"timing"`         // Trigger timing (e.g., "BEFORE", "AFTER", "INSTEAD OF")
	Events         []TriggerEvent     `json:"events"`         // List of events that dispatch the trigger (e.g., "INSERT", "UPDATE", "DELETE")
	Orientation    TriggerOrientation `json:"orientation"`    // Trigger orientation (e.g., "FOR EACH ROW", "FOR EACH STATEMENT")
	Enabled        TriggerEnabled     `json:"enabled"`        // Trigger enabled state (e.g., "ORIGIN", "ALWAYS", "REPLICA", "DISABLED")
	FunctionSchema string             `json:"functionSchema"` // Schema of the associated function (e.g., "public")
	FunctionName   string             `json:"functionName"`   // Name of the associated function called by the trigger
}
