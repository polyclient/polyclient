package database

import (
	"context"
	"database/sql"
)

type SQLQuerier interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}
