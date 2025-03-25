// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/polyclient/polyclient/internal/database"
	"github.com/polyclient/polyclient/internal/datapipe"
	"github.com/urfave/cli/v3"
)

// NewDatabaseCommand returns a new database command for managing databases from the CLI.
func NewDatabaseCommand(driverRegistry *database.Registry[database.Driver]) *cli.Command {
	return &cli.Command{
		Name:  "database",
		Usage: "Manage databases from the CLI",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "driver",
				Aliases:  []string{"d"},
				Usage:    "Database driver name (e.g., 'sqlite', 'postgres', 'mongodb')",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "dsn",
				Usage:    "Data Source Name (DSN) for the database connection",
				Required: true,
			},
		},
		Commands: []*cli.Command{
			newPingCommand(driverRegistry),
			newQueryCommand(driverRegistry),
		},
	}
}

// newPingCommand returns a new ping command that can be used to ping a database from the CLI.
func newPingCommand(driverRegistry *database.Registry[database.Driver]) *cli.Command {
	return &cli.Command{
		Name:  "ping",
		Usage: "Ping a database to check connectivity",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			flagDriver := cmd.String("driver")
			flagDSN := cmd.String("dsn")

			driver, ok := driverRegistry.Get(flagDriver)
			if !ok {
				return fmt.Errorf("driver %s not found", flagDriver)
			}

			switch d := driver.(type) {
			case database.DriverSQL:
				conn, err := d.Connect(flagDSN)
				if err != nil {
					return fmt.Errorf("failed to connect to database: %w", err)
				}

				return conn.PingContext(ctx)

			case database.DriverNoSQL:
				conn, err := d.Connect(flagDSN)
				if err != nil {
					return fmt.Errorf("failed to connect to database: %w", err)
				}

				return conn.PingContext(ctx)

			default:
				return errors.New("driver does not support ping")
			}
		},
	}
}

// newQueryCommand returns a new query command that can be used to query a database from the CLI.
func newQueryCommand(dr *database.Registry[database.Driver]) *cli.Command {
	exportFormats := datapipe.GetAvailableExportFormats()

	return &cli.Command{
		Name:  "query",
		Usage: "Execute a query against a database (SQL or NoSQL)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "query",
				Aliases:  []string{"q"},
				Usage:    "Query to execute (e.g., 'SELECT * FROM users')",
				Required: true,
			},
			&cli.StringFlag{
				Name:        "output",
				Usage:       fmt.Sprintf("Output format (%s)", exportFormats),
				Aliases:     []string{"o"},
				Value:       "json",
				DefaultText: "json",
				Validator: func(output string) error {
					_, ok := datapipe.GetRegistryEntry(datapipe.Format(output))
					if !ok {
						return fmt.Errorf("invalid output format: %s\nSupported formats: %s", output, exportFormats)
					}

					return nil
				},
			},
			&cli.StringFlag{
				Name:        "destination",
				Usage:       "Output destination (stdout or [path-to-file])",
				Aliases:     []string{"d"},
				Value:       "stdout",
				DefaultText: "stdout",
				Validator: func(write string) error {
					if write == "" {
						return errors.New("flag needs an argument: -w")
					}

					return nil
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			flagDriver := cmd.String("driver")
			flagDSN := cmd.String("dsn")
			flagQuery := cmd.String("query")
			flagOutput := cmd.String("output")
			flagDestination := cmd.String("destination")

			var w io.Writer
			if flagDestination == "stdout" {
				w = os.Stdout
			} else {
				file, err := os.Create(path.Clean(flagDestination))
				if err != nil {
					return fmt.Errorf("failed to create file '%s': %w", flagDestination, err)
				}
				defer file.Close()
				w = file
			}

			var resultBytes []byte

			driver, ok := dr.Get(flagDriver)
			if !ok {
				return fmt.Errorf("driver %s not found", flagDriver)
			}

			switch d := driver.(type) {
			case database.DriverSQL:
				conn, err := d.Connect(flagDSN)
				if err != nil {
					return fmt.Errorf("failed to connect to database: %w", err)
				}

				rows, err := conn.QueryContext(ctx, flagQuery)
				if err != nil {
					return fmt.Errorf("failed to execute query: %w", err)
				}
				defer rows.Close()

				columns, err := rows.Columns()
				if err != nil {
					return fmt.Errorf("failed to get columns: %w", err)
				}

				results := []map[string]any{}

				for rows.Next() {
					values := make([]any, len(columns))
					valuesPtrs := make([]any, len(columns))

					for i := range values {
						valuesPtrs[i] = &values[i]
					}

					if err := rows.Scan(valuesPtrs...); err != nil {
						return fmt.Errorf("failed to scan row: %w", err)
					}

					result := make(map[string]any)

					for i, col := range columns {
						result[col] = values[i]
					}

					results = append(results, result)
				}

				if err := rows.Err(); err != nil {
					return fmt.Errorf("failed to iterate rows: %w", err)
				}

				resultBytes, err = json.Marshal(results)
				if err != nil {
					return fmt.Errorf("failed to marshal results: %w", err)
				}

			default:
				return errors.New("driver does not support query")
			}

			result, err := datapipe.ParseDataFromBytes[any](resultBytes, datapipe.Format(flagOutput))
			if err != nil {
				return fmt.Errorf("failed to parse data: %w", err)
			}

			entry, ok := datapipe.GetRegistryEntry(datapipe.Format(flagOutput))
			if !ok {
				return fmt.Errorf("output format %s is not supported", flagOutput)
			}

			if err := entry.Exporter.Export(w, result); err != nil {
				return fmt.Errorf("failed to export data: %w", err)
			}

			return nil
		},
	}
}
