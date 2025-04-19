package postgres

import (
	"context"
	"database/sql"
)

type Connection struct {
	db *sql.DB
}

func (c Connection) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

func (c Connection) Close() error {
	return c.db.Close()
}
