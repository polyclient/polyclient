package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/polyclient/polyclient/internal/database"
)

type Driver struct{}

func NewDriver() database.DriverSQL {
	return &Driver{}
}

// Name implements database.Driver.
func (d *Driver) Name() string {
	return "postgres"
}

// Type implements database.Driver.
func (d *Driver) Type() database.DriverType {
	return database.DriverTypeSQL
}

// CreateConnection implements database.Driver.
func (d *Driver) Connect(dsn string) (database.ConnectionSQL, error) {
	fmt.Println("Creating PostgreSQL connection...")

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	return Connection{db: sqlDB}, nil
}
