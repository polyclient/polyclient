package postgres

import (
	"context"
	"database/sql"

	"github.com/polyclient/polyclient/internal/db"
)

// Connection implements db.Connection.
type Connection struct {
	db *sql.DB
}

// Info implements db.Connection.Info.
func (c *Connection) Info() db.ConnectionInfo {
	return &ConnectionInfo{db: c.db}
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
func (c *ConnectionInfo) ServerVersion() string {
	var version string

	err := c.db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		return "unknown"
	}

	return version
}
