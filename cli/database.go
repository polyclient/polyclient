// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/polyclient/polyclient/drivers/sqlite"
	"github.com/polyclient/polyclient/internal/database"
	"github.com/polyclient/polyclient/internal/datapipe"
	"github.com/urfave/cli/v3"
)

// NewDatabaseCommand returns a new database command for managing databases from the CLI.
func NewDatabaseCommand(dr *database.Registry[database.AnyDriver]) *cli.Command {
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
		},
		Commands: []*cli.Command{
			newPingCommand(dr),
			newQueryCommand(dr),
		},
	}
}

// newPingCommand returns a new ping command that can be used to ping a database from the CLI.
func newPingCommand(dr *database.Registry[database.AnyDriver]) *cli.Command {
	return &cli.Command{
		Name:  "ping",
		Usage: "Ping a database to check connectivity",
		Flags: []cli.Flag{
			// TODO: Make this agnostic
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "Path to the database (for SQLite)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			driverName := cmd.String("driver")
			path := cmd.String("path")

			driver, ok := dr.Get(driverName)
			if !ok {
				return fmt.Errorf("driver %s not found", driverName)
			}

			// TODO: This is just for testing
			config := &sqlite.ConnectionConfig{Path: path}
			conn, err := driver.CreateConnection(config)
			if err != nil {
				return fmt.Errorf("failed to create connection: %w", err)
			}

			if err := conn.PingContext(ctx); err != nil {
				return fmt.Errorf("failed to ping database: %w", err)
			}

			fmt.Printf("Successfully pinged %s database", driverName)
			return nil
		},
	}
}

// newQueryCommand returns a new query command that can be used to query a database from the CLI.
func newQueryCommand(dr *database.Registry[database.AnyDriver]) *cli.Command {
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
		Action: func(_ context.Context, cmd *cli.Command) error {
			// query := cmd.String("query")
			// output := cmd.String("output")
			// destination := cmd.String("destination")

			// var w io.Writer
			// if destination == "stdout" {
			// 	w = os.Stdout
			// } else {
			// 	file, err := os.Create(path.Clean(destination))
			// 	if err != nil {
			// 		return fmt.Errorf("failed to create file '%s': %w", destination, err)
			// 	}
			// 	defer file.Close()
			// 	w = file
			// }

			// resultBytes, err := pr.CallFunction("sqlite", "query", query)
			// if err != nil {
			// 	return fmt.Errorf("failed to execute query: %w", err)
			// }

			// result, err := datapipe.ParseDataFromBytes[any](resultBytes, datapipe.Format(output))
			// if err != nil {
			// 	return fmt.Errorf("failed to parse data: %w", err)
			// }

			// entry, ok := datapipe.GetRegistryEntry(datapipe.Format(output))
			// if !ok {
			// 	return fmt.Errorf("output format %s is not supported", output)
			// }

			// if err := entry.Exporter.Export(w, result); err != nil {
			// 	return fmt.Errorf("failed to export data: %w", err)
			// }

			return nil
		},
	}
}
