package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/polyclient/polyclient/internal/db"
	_ "modernc.org/sqlite"
)

// Driver implements db.Driver.
type Driver struct{}

// Name implements db.Driver.
func (d *Driver) Name() string {
	return "postgres"
}

// Open implements db.Driver.
func (d *Driver) Open(ctx context.Context, config db.Config) (db.Connection, error) {
	dsn, ok := config["dsn"].(string)
	if !ok {
		return nil, errors.New("sqlite: missing or invalid 'dsn' in config")
	}

	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("sqlite: failed to create connection: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		if err := sqlDB.Close(); err != nil {
			return nil, fmt.Errorf("sqlite: failed to close connection: %w", err)
		}

		return nil, fmt.Errorf("sqlite: failed to ping database: %w", err)
	}

	return &Connection{db: sqlDB}, nil
}

// NewDriver returns a new postgres driver.
func NewDriver() db.Driver {
	return &Driver{}
}
