package sqlite

import (
	"context"
	"database/sql"
	"fmt"
)

// ConnectionConfig defines the configuration for SQLite.
type ConnectionConfig struct {
	Path string
}

func (c *ConnectionConfig) Validate() error {
	if c.Path == "" {
		return fmt.Errorf("missing path")
	}

	return nil
}

// Connection implements db.Connection.
type Connection struct {
	db *sql.DB
}

func (c *Connection) Close() error {
	fmt.Println("Closing SQLite connection")
	return c.db.Close()
}

func (c *Connection) Ping() error {
	fmt.Println("Pinging SQLite connection")
	return c.db.Ping()
}

func (c *Connection) PingContext(ctx context.Context) error {
	fmt.Println("Pinging SQLite connection")
	return c.db.PingContext(ctx)
}
