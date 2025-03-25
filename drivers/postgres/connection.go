package postgres

import (
	"context"
	"database/sql"
)

type Connection struct {
	db *sql.DB
}

// Close implements database.Connection.
func (c Connection) Close() error {
	return c.db.Close()
}

// Ping implements database.Connection.
func (c Connection) Ping() error {
	return c.db.Ping()
}

// PingContext implements database.Connection.
func (c Connection) PingContext(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

func (c Connection) Query(query string, args ...any) (*sql.Rows, error) {
	return c.db.Query(query, args...)
}

// Query implements database.ConnectionSQL.
func (c Connection) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return c.db.QueryContext(ctx, query, args...)
}
