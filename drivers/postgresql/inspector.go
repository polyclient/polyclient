package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/polyclient/polyclient/internal/db"
)

const (
	triggerTypeRow      int16 = (1 << 0) // ROW
	triggerTypeBefore   int16 = (1 << 1) // BEFORE
	triggerTypeInsert   int16 = (1 << 2) // INSERT
	triggerTypeDelete   int16 = (1 << 3) // DELETE
	triggerTypeUpdate   int16 = (1 << 4) // UPDATE
	triggerTypeTruncate int16 = (1 << 5) // TRUNCATE
	triggerTypeInstead  int16 = (1 << 6) // INSTEAD
)

var defaultListTablesOptions = func() []db.ListTablesOption {
	return []db.ListTablesOption{
		db.WithTablesSchema("public"),
		db.WithTablesFilter(""),
		db.WithTablesLimit(100),
		db.WithTablesOffset(0),
	}
}

var defaultGetTableOptions = func() []db.GetTableOption {
	return []db.GetTableOption{
		db.WithTableSchema("public"),
	}
}

// ListTables implements db.TableLister.
func (c *Connection) ListTables(ctx context.Context, options ...db.ListTablesOption) ([]db.TableSummary, error) {
	opts := db.NewListTablesOptions(append(defaultListTablesOptions(), options...)...)

	query := `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_schema = $1 AND table_type = 'BASE TABLE' AND table_name LIKE '%%' || $2 || '%%'
		ORDER BY table_name ASC
		LIMIT $3
		OFFSET $4
	`

	rows, err := c.db.QueryContext(ctx, query, opts.Schema, opts.Filter, opts.Limit, opts.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}
	defer rows.Close()

	var summary []db.TableSummary

	for rows.Next() {
		var name string

		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to list tables: %w", err)
		}

		summary = append(summary, db.TableSummary{
			GenericObjectIdentifier: db.GenericObjectIdentifier{
				Name:   name,
				Schema: opts.Schema,
			},
			GenericObjectAttributes: db.GenericObjectAttributes{
				Attributes: nil,
			},
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to list tables: %w", err)
	}

	return summary, nil
}

// GetTable implements db.TableGetter.
func (c *Connection) GetTable(ctx context.Context, name string, options ...db.GetTableOption) (*db.TableDetail, error) {
	opts := db.NewGetTableOptions(append(defaultGetTableOptions(), options...)...)

	var exists bool
	err := c.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables
			WHERE table_schema = $1 AND table_name = $2
		)`, opts.Schema, name).Scan(&exists)

	if err != nil {
		return nil, fmt.Errorf("failed to get table: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("table %s.%s does not exist", opts.Schema, name)
	}

	var tableOID int
	err = c.db.QueryRowContext(ctx, `
		SELECT oid 
        FROM pg_catalog.pg_class 
        WHERE relname = $1 
        AND relnamespace = (SELECT oid 
                            FROM pg_catalog.pg_namespace 
                            WHERE nspname = $2)`,
		name, opts.Schema).Scan(&tableOID)

	if err != nil {
		return nil, fmt.Errorf("failed to get table OID: %w", err)
	}

	detail := db.TableDetail{
		TableSummary: db.TableSummary{
			GenericObjectIdentifier: db.GenericObjectIdentifier{
				Schema: opts.Schema,
				Name:   name,
				OID:    tableOID,
			},
		},
	}

	detail.Columns, err = c.fetchTableColumns(ctx, opts.Schema, name, tableOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table columns: %w", err)
	}

	detail.Indexes, err = c.fetchTableIndexes(ctx, opts.Schema, tableOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table indexes: %w", err)
	}

	detail.Constraints, err = c.fetchTableConstraints(ctx, opts.Schema, tableOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table constraints: %w", err)
	}

	detail.Triggers, err = c.fetchTableTriggers(ctx, opts.Schema, tableOID)
	if err != nil {
		return nil, fmt.Errorf("failed to get table triggers: %w", err)
	}

	return &detail, nil
}

// fetchTableColumns retrieves column details for a table.
func (c *Connection) fetchTableColumns(ctx context.Context, schema, tableName string, tableOID int) ([]db.ColumnDetail, error) {
	query := `
		SELECT
			c.column_name,
			c.udt_name,
			c.data_type,
			c.character_maximum_length,
			c.numeric_precision,
			c.numeric_scale,
			c.ordinal_position,
			CASE WHEN c.is_nullable = 'YES' THEN TRUE ELSE FALSE END AS is_nullable,
			CASE WHEN c.is_self_referencing = 'YES' THEN TRUE ELSE FALSE END AS is_self_referencing,
			CASE WHEN c.is_updatable = 'YES' THEN TRUE ELSE FALSE END AS is_updatable,
			CASE WHEN c.is_generated = 'ALWAYS' THEN TRUE ELSE FALSE END AS is_generated,
			c.generation_expression,
			c.column_default,
			pg_catalog.col_description($3, c.ordinal_position) AS comment
		FROM information_schema.columns c
		WHERE c.table_schema = $2 AND c.table_name = $1
		ORDER BY c.ordinal_position;
	`

	rows, err := c.db.QueryContext(ctx, query, tableName, schema, tableOID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch columns: %w", err)
	}
	defer rows.Close()

	var columns []db.ColumnDetail

	for rows.Next() {
		var columnName, typeName, typeFamily string
		var ordinalPosition int
		var maxLength, precision, scale sql.NullInt64
		var isNullable, isSelfReferenced, isUpdatable, isGenerated bool
		var generationExpr, defaultValue, comment sql.NullString

		if err := rows.Scan(
			&columnName,
			&typeName,
			&typeFamily,
			&maxLength,
			&precision,
			&scale,
			&ordinalPosition,
			&isNullable,
			&isSelfReferenced,
			&isUpdatable,
			&isGenerated,
			&generationExpr,
			&defaultValue,
			&comment,
		); err != nil {
			return nil, fmt.Errorf("failed to scan column: %w", err)
		}

		column := db.ColumnDetail{
			ColumnSummary: db.ColumnSummary{
				GenericObjectIdentifier: db.GenericObjectIdentifier{
					Name:   columnName,
					Schema: schema,
				},
				GenericObjectAttributes: db.GenericObjectAttributes{
					Attributes: nil,
				},
			},
			GenericObjectComment: db.GenericObjectComment{
				Comment: comment.String,
			},
			TypeName:         typeName,
			TypeFamily:       typeFamily,
			OrdinalPosition:  ordinalPosition,
			MaxLength:        int(maxLength.Int64),
			Precision:        int(precision.Int64),
			Scale:            int(scale.Int64),
			IsNullable:       isNullable,
			IsSelfReferenced: isSelfReferenced,
			IsUpdatable:      isUpdatable,
			IsGenerated:      isGenerated,
			GenerationExpr:   generationExpr.String,
			DefaultValue:     defaultValue.String,
		}

		columns = append(columns, column)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate columns: %w", err)
	}

	return columns, nil
}

// fetchTableIndexes retrieves index details for a table.
func (c *Connection) fetchTableIndexes(ctx context.Context, schema string, tableOID int) ([]db.IndexDetail, error) {
	query := `
		SELECT
			i.indexrelid,
			c.relname AS index_name,
			pg_catalog.pg_get_indexdef(i.indexrelid) AS definition,
			pg_catalog.obj_description(i.indexrelid, 'pg_class') AS comment,
			i.indisunique,
			i.indisprimary,
			(
				SELECT array_agg(a.attname ORDER BY u.ord)
				FROM (
					SELECT attnum, row_number() OVER () AS ord
					FROM unnest(i.indkey) AS attnum
				) u
				JOIN pg_attribute a ON a.attrelid = i.indrelid AND a.attnum = u.attnum
			) AS column_names
		FROM pg_catalog.pg_index i
		JOIN pg_catalog.pg_class c ON c.oid = i.indexrelid
		WHERE i.indrelid = $1
		ORDER BY c.relname;
	`

	rows, err := c.db.QueryContext(ctx, query, tableOID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch indexes: %w", err)
	}
	defer rows.Close()

	var indexes []db.IndexDetail

	for rows.Next() {
		var indexOID int
		var indexName, definition, comment sql.NullString
		var isUnique, isPrimary bool
		var columnNames string

		if err := rows.Scan(
			&indexOID,
			&indexName,
			&definition,
			&comment,
			&isUnique,
			&isPrimary,
			&columnNames,
		); err != nil {
			fmt.Printf("Index scan error: %v\n", err)
			return nil, fmt.Errorf("failed to scan index: %w", err)
		}

		indexes = append(indexes, db.IndexDetail{
			IndexSummary: db.IndexSummary{
				GenericObjectIdentifier: db.GenericObjectIdentifier{
					Name:   indexName.String,
					Schema: schema,
					OID:    indexOID,
				},
			},
			GenericObjectComment: db.GenericObjectComment{
				Comment: comment.String,
			},
			GenericObjectDefinition: db.GenericObjectDefinition{
				Definition: definition.String,
			},
			Columns:    columnNames,
			Definition: definition.String,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate indexes: %w", err)
	}

	return indexes, nil
}

// fetchTableConstraints retrieves constraint details for a table.
func (c *Connection) fetchTableConstraints(ctx context.Context, schema string, tableOID int) ([]db.ConstraintDetail, error) {
	query := `
		SELECT
			c.oid,
			c.conname AS constraint_name,
			c.contype,
			pg_catalog.pg_get_constraintdef(c.oid) AS definition,
			pg_catalog.obj_description(c.oid, 'pg_constraint') AS comment,    
			f.relnamespace::regnamespace::text AS foreign_schema,
    		f.relname AS foreign_table,
			c.confupdtype AS on_update_action,
			c.confdeltype AS on_delete_action,
			c.confmatchtype AS match_option
		FROM pg_catalog.pg_constraint c
		LEFT JOIN pg_catalog.pg_class f ON f.oid = c.confrelid
		WHERE c.conrelid = $1
		ORDER BY c.contype;
	`

	rows, err := c.db.QueryContext(ctx, query, tableOID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch constraints: %w", err)
	}
	defer rows.Close()

	var constraints []db.ConstraintDetail

	for rows.Next() {
		var constraintOID int
		var constraintName, constraintType, definition, comment, foreignSchema, foreignTable sql.NullString
		var onUpdateAction, onDeleteAction, matchOption sql.NullString

		if err := rows.Scan(
			&constraintOID,
			&constraintName,
			&constraintType,
			&definition,
			&comment,
			&foreignSchema,
			&foreignTable,
			&onUpdateAction,
			&onDeleteAction,
			&matchOption,
		); err != nil {
			return nil, fmt.Errorf("failed to scan constraint: %w", err)
		}

		constraints = append(constraints, db.ConstraintDetail{
			ConstraintSummary: db.ConstraintSummary{
				GenericObjectIdentifier: db.GenericObjectIdentifier{
					Name:   constraintName.String,
					Schema: schema,
					OID:    constraintOID,
				},
				Type: mapConstraintTypeToCode(constraintType.String),
			},
			GenericObjectComment: db.GenericObjectComment{
				Comment: comment.String,
			},
			GenericObjectDefinition: db.GenericObjectDefinition{
				Definition: definition.String,
			},
			ForeignSchema:  foreignSchema.String,
			ForeignTable:   foreignTable.String,
			OnUpdateAction: mapConstraintActionToCode(onUpdateAction.String),
			OnDeleteAction: mapConstraintActionToCode(onDeleteAction.String),
			MatchOption:    mapConstraintMatchToCode(matchOption.String),
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate constraints: %w", err)
	}

	return constraints, nil
}

// fetchTableTriggers retrieves trigger details for a table.
func (c *Connection) fetchTableTriggers(ctx context.Context, schema string, tableOID int) ([]db.TriggerDetail, error) {
	query := `
		SELECT
			tg.oid,
			tg.tgname AS trigger_name,
			pg_catalog.pg_get_triggerdef(tg.oid) AS definition,
			tg.tgtype AS type_flags, -- Bitmask for timing, event, orientation
			tg.tgenabled AS enabled_char, -- 'O', 'D', 'R', 'A'
			p.proname AS function_name,
			pn.nspname AS function_schema,
			pg_catalog.obj_description(tg.oid, 'pg_trigger') AS comment
		FROM pg_catalog.pg_trigger tg
		JOIN pg_catalog.pg_proc p ON p.oid = tg.tgfoid
		JOIN pg_catalog.pg_namespace pn ON pn.oid = p.pronamespace
		WHERE tg.tgrelid = $1 AND tg.tgisinternal = false -- Exclude internal triggers
		ORDER BY tg.tgname;
	`

	rows, err := c.db.QueryContext(ctx, query, tableOID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch triggers: %w", err)
	}
	defer rows.Close()

	var triggers []db.TriggerDetail

	for rows.Next() {
		var triggerOID int
		var triggerName, definition, enabledChar, functionName, functionSchema, comment sql.NullString
		var typeFlags int16

		if err := rows.Scan(
			&triggerOID,
			&triggerName,
			&definition,
			&typeFlags,
			&enabledChar,
			&functionName,
			&functionSchema,
			&comment,
		); err != nil {
			return nil, fmt.Errorf("failed to scan trigger: %w", err)
		}

		enabledStatus := mapTriggerEnabledToCode(enabledChar.String)
		timing, events, orientation := mapTriggerTypeFlags(typeFlags)

		triggers = append(triggers, db.TriggerDetail{
			TriggerSummary: db.TriggerSummary{
				GenericObjectIdentifier: db.GenericObjectIdentifier{
					Name:   triggerName.String,
					Schema: schema,
					OID:    triggerOID,
				},
			},
			GenericObjectComment: db.GenericObjectComment{
				Comment: comment.String,
			},
			GenericObjectDefinition: db.GenericObjectDefinition{
				Definition: definition.String,
			},
			Timing:         timing,
			Events:         events,
			Orientation:    orientation,
			Enabled:        enabledStatus,
			FunctionSchema: functionSchema.String,
			FunctionName:   functionName.String,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate triggers: %w", err)
	}

	return triggers, nil
}

// mapConstraintTypeToCode maps a constraint type code to a ConstraintType value.
func mapConstraintTypeToCode(code string) db.ConstraintType {
	switch code {
	case "p":
		return db.ConstraintTypePrimaryKey
	case "f":
		return db.ConstraintTypeForeignKey
	case "u":
		return db.ConstraintTypeUnique
	case "c":
		return db.ConstraintTypeCheck
	default:
		return ""
	}
}

// mapConstraintActionToCode maps a constraint action code to a ConstraintActionType value.
func mapConstraintActionToCode(code string) db.ConstraintActionType {
	switch code {
	case "a":
		return db.ConstraintActionTypeNoAction
	case "r":
		return db.ConstraintActionTypeRestrict
	case "c":
		return db.ConstraintActionTypeCascade
	case "n":
		return db.ConstraintActionTypeSetNull
	case "d":
		return db.ConstraintActionTypeSetDefault
	default:
		return ""
	}
}

// mapConstraintMatchToCode maps a constraint match code to a ConstraintMatchType value.
func mapConstraintMatchToCode(code string) db.ConstraintMatchType {
	switch code {
	case "f":
		return db.ConstraintMatchTypeFull
	case "p":
		return db.ConstraintMatchTypePartial
	case "s":
		return db.ConstraintMatchTypeSimple
	default:
		return ""
	}
}

// mapTriggerEnabledToCode maps a trigger enabled code to a TriggerEnabled value.
func mapTriggerEnabledToCode(code string) db.TriggerEnabled {
	switch code {
	case "O":
		return db.TriggerEnabledOrigin
	case "D":
		return db.TriggerEnabledDisabled
	case "R":
		return db.TriggerEnabledReplica
	case "A":
		return db.TriggerEnabledAlways
	default:
		return ""
	}
}

// mapTriggerTypeFlags decodes the pg_trigger.tgtype bitmask.
// Reference: https://www.postgresql.org/docs/current/catalog-pg-trigger.html
func mapTriggerTypeFlags(flags int16) (db.TriggerTiming, []db.TriggerEvent, db.TriggerOrientation) {
	var timing db.TriggerTiming
	var events []db.TriggerEvent
	var orientation db.TriggerOrientation

	// Orientation (Bit 0)
	if (flags & triggerTypeRow) != 0 {
		orientation = db.TriggerOrientationRow
	} else {
		orientation = db.TriggerOrientationStatement
	}

	// Timing (Bits 1 and 6)
	isTruncate := (flags & triggerTypeTruncate) != 0
	if isTruncate {
		// TRUNCATE triggers are always AFTER and STATEMENT
		timing = db.TriggerTimingAfter
		orientation = db.TriggerOrientationStatement
	} else if (flags & triggerTypeInstead) != 0 {
		timing = db.TriggerTimingInsteadOf
	} else if (flags & triggerTypeBefore) != 0 {
		timing = db.TriggerTimingBefore
	} else {
		timing = db.TriggerTimingAfter
	}

	// Events (Bits 2, 3, 4, 5)
	if (flags & triggerTypeInsert) != 0 {
		events = append(events, db.TriggerEventInsert)
	}

	if (flags & triggerTypeDelete) != 0 {
		events = append(events, db.TriggerEventDelete)
	}

	if (flags & triggerTypeUpdate) != 0 {
		events = append(events, db.TriggerEventUpdate)
	}

	if isTruncate {
		// Add TRUNCATE event if the flag was set
		events = append(events, db.TriggerEventTruncate)
	}

	return timing, events, orientation
}

var _ db.TableLister = (*Connection)(nil)
var _ db.TableGetter = (*Connection)(nil)
