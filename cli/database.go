// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/polyclient/polyclient/internal/datapipe"
	"github.com/polyclient/polyclient/internal/db"
	"github.com/polyclient/polyclient/internal/sdk"
	"github.com/urfave/cli/v3"
)

// NewDatabaseCommand returns a new database command for managing databases and their connections.
func NewDatabaseCommand(dbSDK *sdk.DatabaseSDK) *cli.Command {
	return &cli.Command{
		Name:  "database",
		Usage: "Manage databases and their connections",
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
			newPingCommand(dbSDK),
			newListTablesCommand(dbSDK),
			newQueryCommand(dbSDK),
		},
	}
}

// newPingCommand returns a new ping command that can be used to ping a database from the CLI.
func newPingCommand(dbSDK *sdk.DatabaseSDK) *cli.Command {
	return &cli.Command{
		Name:  "ping",
		Usage: "Ping a database to check connectivity",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			flagDriver := cmd.String("driver")
			flagDSN := cmd.String("dsn")

			conn, err := dbSDK.OpenConnection(ctx, flagDriver, db.Config{"dsn": flagDSN})
			if err != nil {
				return fmt.Errorf("failed to open connection: %w", err)
			}

			if err := conn.Ping(ctx); err != nil {
				return fmt.Errorf("failed to ping database: %w", err)
			}

			info := conn.Info()
			fmt.Printf("Connected to %s", info.ServerVersion())

			return nil
		},
	}
}

func newListTablesCommand(dbSDK *sdk.DatabaseSDK) *cli.Command {
	return &cli.Command{
		Name:  "list-tables",
		Usage: "List all tables in a database",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			flagDriver := cmd.String("driver")
			flagDSN := cmd.String("dsn")

			conn, err := dbSDK.OpenConnection(ctx, flagDriver, db.Config{"dsn": flagDSN})
			if err != nil {
				return fmt.Errorf("failed to open connection: %w", err)
			}

			tables, err := conn.Schema().ListTables(ctx)
			if err != nil {
				return fmt.Errorf("failed to list tables: %w", err)
			}

			for _, table := range tables {
				fmt.Println(table.Name)
			}

			return nil
		},
	}
}

// newQueryCommand returns a new query command that can be used to query a database from the CLI.
func newQueryCommand(dbSDK *sdk.DatabaseSDK) *cli.Command {
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
						return fmt.Errorf(
							"invalid output format: %s\nSupported formats: %s",
							output,
							exportFormats,
						)
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

			conn, err := dbSDK.OpenConnection(ctx, flagDriver, db.Config{"dsn": flagDSN})
			if err != nil {
				return fmt.Errorf("failed to open connection: %w", err)
			}

			result, err := conn.Query().Execute(ctx, flagQuery)
			if err != nil {
				return fmt.Errorf("failed to execute query: %w", err)
			}

			fmt.Println(result) // TODO: connect this to datapipe

			entry, ok := datapipe.GetRegistryEntry(datapipe.Format(flagOutput))
			if !ok {
				return fmt.Errorf("output format %s is not supported", flagOutput)
			}

			if err := entry.Export(w, result); err != nil {
				return fmt.Errorf("failed to export data: %w", err)
			}

			return nil
		},
	}
}
