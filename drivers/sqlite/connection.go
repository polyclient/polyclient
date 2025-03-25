package sqlite

import (
	"context"
	"database/sql"
)

// Connection implements db.Connection.
type Connection struct {
	db *sql.DB
}

func (c *Connection) Close() error {
	return c.db.Close()
}

func (c *Connection) Ping() error {
	return c.db.Ping()
}

func (c *Connection) PingContext(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

func (c *Connection) Query(query string, args ...any) (*sql.Rows, error) {
	return c.db.Query(query, args...)
}

func (c *Connection) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return c.db.QueryContext(ctx, query, args...)
}
