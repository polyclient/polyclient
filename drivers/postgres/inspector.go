package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/polyclient/polyclient/internal/db"
)

// ListTables implements db.TableLister.
func (c *Connection) ListTables(ctx context.Context) ([]db.TableSummary, error) {
	query := `	
		SELECT table_name
		FROM information_schema.tables
		WHERE table_type = 'public'
		ORDER BY table_name
	`

	rows, err := c.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to list tables: %w", err)
	}
	defer rows.Close()

	var summary []db.TableSummary

	for rows.Next() {
		var name string

		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("postgres: failed to list tables: %w", err)
		}

		summary = append(summary, db.TableSummary{
			GenericObjectIdentifier: db.GenericObjectIdentifier{
				Name:   name,
				Schema: "public",
			},
			GenericObjectOwner: db.GenericObjectOwner{
				Owner: "postgres", // TODO: get owner
			},
			GenericObjectAttributes: db.GenericObjectAttributes{
				Attributes: nil,
			},
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres: failed to list tables: %w", err)
	}

	return summary, nil
}

// GetTable implements db.TableGetter.
func (c *Connection) GetTable(ctx context.Context, name string) (db.TableDetail, error) {
	query := `	
		SELECT table_name, comment, create_time
		FROM information_schema.tables
		WHERE table_type = 'public' AND table_name = $1
		ORDER BY ordinal_position
	`

	rows, err := c.db.QueryContext(ctx, query, name)
	if err != nil {
		return db.TableDetail{}, fmt.Errorf("postgres: failed to get table: %w", err)
	}
	defer rows.Close()

	detail := db.TableDetail{
		TableSummary: db.TableSummary{
			GenericObjectIdentifier: db.GenericObjectIdentifier{
				Schema: "public",
				Name:   name,
			},
		},
	}

	for rows.Next() {
		var name string
		var dataType string
		var isNullable bool
		var comment string
		var createdAt time.Time

		if err := rows.Scan(&name, &dataType, &isNullable, &comment, &createdAt); err != nil {
			return db.TableDetail{}, fmt.Errorf("postgres: failed to get table: %w", err)
		}

		detail.Columns = append(detail.Columns, db.ColumnDetail{
			ColumnSummary: db.ColumnSummary{
				GenericObjectIdentifier: db.GenericObjectIdentifier{
					Name:   name,
					Schema: "public",
				},
				GenericObjectOwner: db.GenericObjectOwner{
					Owner: "postgres", // TODO: get owner
				},
				GenericObjectAttributes: db.GenericObjectAttributes{
					Attributes: nil,
				},
			},
			GenericObjectComment: db.GenericObjectComment{
				Comment: comment,
			},
			Position:     0, // TODO: get position
			BaseType:     dataType,
			FullDataType: dataType,
			IsNullable:   isNullable,
			IsPrimaryKey: false,
			DefaultValue: nil,
		})
	}

	if err := rows.Err(); err != nil {
		return db.TableDetail{}, fmt.Errorf("postgres: failed to get table: %w", err)
	}

	return detail, nil
}

var _ db.TableLister = (*Connection)(nil)
var _ db.TableGetter = (*Connection)(nil)
