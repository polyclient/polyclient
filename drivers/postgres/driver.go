package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/polyclient/polyclient/internal/db"
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
		return nil, errors.New("postgres: missing or invalid 'dsn' in config")
	}

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres: failed to create connection: %w", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		if err := sqlDB.Close(); err != nil {
			return nil, fmt.Errorf("postgres: failed to close connection: %w", err)
		}

		return nil, fmt.Errorf("postgres: failed to ping database: %w", err)
	}

	return &Connection{db: sqlDB}, nil
}

// NewDriver returns a new postgres driver.
func NewDriver() db.Driver {
	return &Driver{}
}
