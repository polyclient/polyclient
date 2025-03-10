// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package shell

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/samber/lo"

	"github.com/polyclient/polyclient/internal/runtime/plugin"
	"github.com/polyclient/polyclient/pkg/dataexchange"
	"github.com/urfave/cli/v3"
)

// NewDatabaseCommand returns a new database command for managing databases from the CLI.
func NewDatabaseCommand(pr *plugin.PluginRegistry) *cli.Command {
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
			newQueryCommand(pr),
		},
	}
}

// newQueryCommand returns a new query command that can be used to query a database from the CLI.
func newQueryCommand(pr *plugin.PluginRegistry) *cli.Command {
	supportedFormats := lo.Map(dataexchange.GetSupportedExportFormats(), func(format dataexchange.Format, _ int) string {
		return fmt.Sprintf("%v", format)
	})

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
				Value:       "json",
				DefaultText: "json",
				Validator: func(output string) error {
					_, ok := dataexchange.GetRegistryEntry(dataexchange.Format(output))
					if !ok {
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
			query := cmd.String("query")
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

			resultBytes, err := pr.CallFunction("sqlite", "query", []byte(query))
			if err != nil {
				return fmt.Errorf("failed to execute query: %w", err)
			}

			result, err := dataexchange.ParseDataFromBytes[any](resultBytes, dataexchange.Format(output))
			if err != nil {
				return fmt.Errorf("failed to parse data: %w", err)
			}

			entry, ok := dataexchange.GetRegistryEntry(dataexchange.Format(output))
			if !ok {
				return fmt.Errorf("output format %s is not supported", output)
			}

			if err := entry.Exporter.Export(w, result); err != nil {
				return fmt.Errorf("failed to export data: %w", err)
			}

			return nil
		},
	}
}
