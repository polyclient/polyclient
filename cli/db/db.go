// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package db

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/polyclient/polyclient/internal/datapipe"
	"github.com/polyclient/polyclient/internal/engine"
	"github.com/urfave/cli/v3"
)

// NewDatabaseCommand returns a new database command for managing databases and their connections.
func NewDatabaseCommand(e *engine.Engine) *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "Manage databases and their connections",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "connection",
				Aliases: []string{"c"},
				Usage:   "Database connection name",
				Config: cli.StringConfig{
					TrimSpace: true,
				},
				Validator: func(connection string) error {
					profile, err := e.SDK.GetManager().GetStore().GetProfile(context.Background(), connection)
					if err != nil || profile == nil {
						return errors.New("connection not found; use 'polyclient db connection list' to see available connections")
					}

					return nil
				},
			},
		},
		Commands: []*cli.Command{
			newConnectionCommand(e),
			newTableCommand(e),
		},
	}
}

func getConnectionName(ctx context.Context, cmd *cli.Command) (string, error) {
	flagConnection := cmd.String("connection")

	if flagConnection == "" {
		return "", errors.New("connection name not set: specify it with --connection or -c (e.g., polyclient db --connection <name>)")
	}

	return flagConnection, nil
}

// newQueryCommand returns a new query command that can be used to query a database from the CLI.
func newQueryCommand(e *engine.Engine) *cli.Command {
	exportFormats := datapipe.GetAvailableExportFormats()

	return &cli.Command{
		Name:  "query",
		Usage: "Execute a query against a database (SQL or NoSQL)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "connection",
				Aliases:  []string{"c"},
				Usage:    "Database connection name",
				Required: true,
			},
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
			flagConnection := cmd.String("connection")
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

			result, err := e.SDK.Query().Execute(ctx, flagConnection, flagQuery)
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
