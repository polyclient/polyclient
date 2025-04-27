package settings

import "github.com/go-playground/validator/v10"

// Drivers holds the settings for the database drivers.
type Drivers struct {
	SQLite     SQLite     `json:"sqlite" validate:"required,dive"`
	PostgreSQL PostgreSQL `json:"postgresql" validate:"required,dive"`
}

// SQLite holds the settings for the SQLite database driver.
type SQLite struct {
	Enabled bool `json:"enabled" validate:"boolean"`
}

// PostgreSQL holds the settings for the PostgreSQL database driver.
type PostgreSQL struct {
	Enabled bool `json:"enabled" validate:"boolean"`
}

// DriversOption is a function that configures the settings for the database drivers.
type DriversOption func(*Drivers)

// WithSQLite sets the settings for the SQLite database driver.
func WithSQLite(enabled bool) DriversOption {
	return func(d *Drivers) {
		d.SQLite.Enabled = enabled
	}
}

// WithPostgreSQL sets the settings for the PostgreSQL database driver.
func WithPostgreSQL(enabled bool) DriversOption {
	return func(d *Drivers) {
		d.PostgreSQL.Enabled = enabled
	}
}

// defaultDrivers returns a new Drivers instance with the default values.
func defaultDriversSettings() *Drivers {
	return &Drivers{
		SQLite:     SQLite{Enabled: true},
		PostgreSQL: PostgreSQL{Enabled: true},
	}
}

// NewDrivers returns a new Drivers instance with the default values.
func NewDrivers(opts ...DriversOption) *Drivers {
	settings := defaultDriversSettings()

	for _, opt := range opts {
		opt(settings)
	}

	return settings
}

// Validate validates the drivers settings.
func (d *Drivers) Validate() error {
	return validator.New(validator.WithRequiredStructEnabled()).Struct(d)
}
