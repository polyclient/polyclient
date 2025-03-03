// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/polyclient/polyclient/internal/exporter"
	"github.com/urfave/cli/v3"
)

// NewDatabaseCommand returns a new database command for managing databases from the CLI.
func NewDatabaseCommand() *cli.Command {
	return &cli.Command{
		Name:  "database",
		Usage: "Manage databases from the CLI",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "connection",
				Aliases:  []string{"c"},
				Usage:    "Connection ID to use",
				Required: true,
			},
		},
		Commands: []*cli.Command{
			newQueryCommand(),
		},
	}
}

// newQueryCommand returns a new query command that can be used to query a database from the CLI.
func newQueryCommand() *cli.Command {
	supportedFormats := strings.Join(exporter.GetSupportedFormats(), ", ")

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
				Usage:       fmt.Sprintf("Output format (%s)", supportedFormats),
				Aliases:     []string{"o"},
				Value:       "markdown",
				DefaultText: "markdown",
				Validator: func(output string) error {
					if !exporter.ValidateFormat(exporter.Format(output)) {
						return fmt.Errorf("invalid output format: %s\nSupported formats: %s", output, supportedFormats)
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
			// query := cmd.String("query")
			output := cmd.String("output")
			destination := cmd.String("destination")

			var w io.Writer
			if destination == "stdout" {
				w = os.Stdout
			} else {
				file, err := os.Create(destination)
				if err != nil {
					return fmt.Errorf("failed to create file '%s': %w", destination, err)
				}
				defer file.Close()
				w = file
			}

			dataExporter := exporter.NewDataExporter(exporter.DataExporterOptions{
				Format: output,
				Writer: w,
			})

			mockData := []struct {
				Name     string
				Email    string
				Birthday time.Time
			}{
				{
					Name:     "John Doe",
					Email:    "johndoe@example.com",
					Birthday: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Name:     "Jane Doe",
					Email:    "janedoe@example.com",
					Birthday: time.Date(1995, 1, 1, 0, 0, 0, 0, time.UTC),
				},
				{
					Name:     "Bob Smith",
					Email:    "bobsmith@example.com",
					Birthday: time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			}

			return dataExporter.Export(mockData)
		},
	}
}
