package postgresql

import (
	"context"
	"database/sql"

	"github.com/polyclient/polyclient/internal/db"
)

// Connection implements db.Connection.
type Connection struct {
	db *sql.DB
}

// Ping implements db.Connection.Ping.
func (c Connection) Ping(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

// Close implements db.Connection.Close.
func (c Connection) Close() error {
	return c.db.Close()
}

// ConnectionInfo implements db.ConnectionInfo.
type ConnectionInfo struct {
	db *sql.DB
}

// ServerVersion implements db.ConnectionInfo.ServerVersion.
func (c *ConnectionInfo) ServerVersion(ctx context.Context) string {
	var version string

	err := c.db.QueryRowContext(ctx, "SELECT version()").Scan(&version)
	if err != nil {
		return "unknown"
	}

	return version
}

// CurrentDatabase implements db.ConnectionInfo.CurrentDatabase.
func (c *ConnectionInfo) CurrentDatabase(ctx context.Context) string {
	var database string

	err := c.db.QueryRowContext(ctx, "SELECT current_database()").Scan(&database)
	if err != nil {
		return "unknown"
	}

	return database
}

var _ db.Connection = (*Connection)(nil)
var _ db.ConnectionInfo = (*ConnectionInfo)(nil)
