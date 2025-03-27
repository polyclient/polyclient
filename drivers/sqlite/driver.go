package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/polyclient/polyclient/internal/database"
	_ "modernc.org/sqlite"
)

type Driver struct{}

func NewDriver() database.DriverSQL {
	return &Driver{}
}

func (d *Driver) Name() string {
	return "sqlite"
}

func (d *Driver) Type() database.DriverType {
	return database.DriverTypeSQL
}

func (d *Driver) Connect(dsn string) (database.ConnectionSQL, error) {
	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	return &Connection{db: sqlDB}, nil
}
