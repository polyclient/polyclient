// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package settings

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema"
)

const schemaFilePath = "jsonschema/settings.schema.json"

// Settings holds the settings for PolyClient.
type Settings struct {
	Drivers Drivers `json:"drivers"`
	API     API     `json:"api"`
	GUI     GUI     `json:"gui"`
	Logging Logging `json:"logging"`
}

// Drivers holds the settings for the database drivers.
type Drivers struct {
	SQLite     SQLite     `json:"sqlite"`
	PostgreSQL PostgreSQL `json:"postgresql"`
}

// SQLite holds the settings for the SQLite database driver.
type SQLite struct {
	Enabled bool `json:"enabled"`
}

// PostgreSQL holds the settings for the PostgreSQL database driver.
type PostgreSQL struct {
	Enabled bool `json:"enabled"`
}

// API holds the settings for the API.
type API struct {
	Host        string         `json:"host"`
	Port        int            `json:"port"`
	CORS        APICORS        `json:"cors"`
	RateLimit   APIRateLimit   `json:"rateLimit"`
	Compression APICompression `json:"compression"`
	Timeouts    APITimeouts    `json:"timeouts"`
}

// APICORS configures Cross-Origin Resource Sharing.
type APICORS struct {
	Enabled        bool     `json:"enabled"`
	AllowedOrigins []string `json:"allowedOrigins"`
	AllowedMethods []string `json:"allowedMethods"`
	AllowedHeaders []string `json:"allowedHeaders"`
	MaxAge         int      `json:"maxAge"`
}

// APIRateLimit configures rate limiting.
type APIRateLimit struct {
	Enabled             bool `json:"enabled"`
	RequestsPerMinute   int  `json:"requestsPerMinute"`
	WindowLengthSeconds int  `json:"windowLengthSeconds"`
}

// APICompression configures compression.
type APICompression struct {
	Enabled bool `json:"enabled"`
	Level   int  `json:"level"`
}

// APICache configures caching.
type APICache struct {
	Enabled           bool    `json:"enabled"`
	DefaultTTLSeconds int     `json:"defaultTTLSeconds"`
	MaxSizeMB         float64 `json:"maxSizeMB"`
}

// APITimeouts configures request timeouts.
type APITimeouts struct {
	Enabled      bool `json:"enabled"`
	ReadSeconds  int  `json:"readSeconds"`
	WriteSeconds int  `json:"writeSeconds"`
	IdleSeconds  int  `json:"idleSeconds"`
}

// GUI holds the settings for the GUI.
type GUI struct {
	Enabled bool   `json:"enabled"`
	Path    string `json:"path"`
	Theme   string `json:"theme"`
	Locale  string `json:"locale"`
}

// Logging holds the settings for logging.
type Logging struct {
	Format   string          `json:"format"`
	Level    string          `json:"level"`
	Rotation LoggingRotation `json:"rotation"`
}

// LoggingRotation holds the settings for log rotation.
type LoggingRotation struct {
	Enabled    bool `json:"enabled"`
	MaxSizeMB  int  `json:"maxSizeMB"`
	MaxAgeDays int  `json:"maxAgeDays"`
	MaxBackups int  `json:"maxBackups"`
	LocalTime  bool `json:"localTime"`
	Compress   bool `json:"compress"`
}

// Validate validates the settings.
func (s *Settings) Validate() error {
	c := jsonschema.NewCompiler()

	schema, err := c.Compile(schemaFilePath)
	if err != nil {
		return fmt.Errorf("failed to compile schema: %w", err)
	}

	settingsJSON, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to marshal settings to JSON: %w", err)
	}

	if err := schema.Validate(bytes.NewReader(settingsJSON)); err != nil {
		return fmt.Errorf("failed to validate settings: %w", err)
	}

	return nil
}
