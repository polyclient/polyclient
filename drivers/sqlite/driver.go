package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/polyclient/polyclient/internal/database"
	_ "modernc.org/sqlite"
)

type Driver struct{}

func NewDriver() database.Driver[database.Connection] {
	return &Driver{}
}

func (d *Driver) Name() string {
	return "sqlite"
}

func (d *Driver) Type() database.DriverType {
	return database.DriverTypeSQL
}

func (d *Driver) CreateConnection(config database.ConnectionConfig) (database.Connection, error) {
	fmt.Println("Creating SQLite connection...")

	cfg, ok := config.(*ConnectionConfig)
	if !ok {
		return nil, fmt.Errorf("invalid connection config")
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	sqlDB, err := sql.Open(d.Name(), cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	return &Connection{db: sqlDB}, nil
}
