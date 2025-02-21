// Copyright (C) 2025 Juan Mesa and contributors
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3
// as published by the Free Software Foundation, with the Runtime
// Library Exception. See the COPYING.RUNTIME file for details.

package db

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/polyclient/polyclient/internal/exporter"
	"github.com/urfave/cli/v3"
)

// NewQueryCommand returns a new query command that can be used to query a database from the CLI
func NewQueryCommand() *cli.Command {
	return &cli.Command{
		Name:     "query",
		Usage:    "Execute a query against a database",
		Category: "Database",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "query",
				Aliases:  []string{"q"},
				Usage:    "Query to execute (e.g. 'SELECT * FROM users')",
				Required: true,
			},
			&cli.StringFlag{
				Name:        "output",
				Usage:       fmt.Sprintf("Output format (%s)", strings.Join(exporter.GetSupportedFormats(), ", ")),
				Aliases:     []string{"o"},
				Value:       "markdown",
				DefaultText: "markdown",
				Validator: func(output string) error {
					if !exporter.ValidateFormat(exporter.Format(output)) {
						return fmt.Errorf("invalid output format: %s", output)
					}

					return nil
				},
			},
			&cli.StringFlag{
				Name:        "write",
				Usage:       "Write method (stdout or file)",
				Aliases:     []string{"w"},
				Value:       "stdout",
				DefaultText: "stdout",
				Validator: func(output string) error {
					if output != "stdout" && output != "file" {
						return fmt.Errorf("invalid write method: %s", output)
					}

					return nil
				},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			query := cmd.String("query")
			output := cmd.String("output")
			write := cmd.String("write")

			fmt.Printf("Query: %s\n", query)
			fmt.Printf("Output: %s\n", output)
			fmt.Printf("Write: %s\n", write)

			var w io.Writer
			if write == "stdout" {
				w = os.Stdout
			} else {
				w = io.Discard
			}

			dataExporter := exporter.NewDataExporter(exporter.DataExporterOptions{
				Format: output,
				Output: w,
			})

			mockData := []map[string]string{
				{
					"foo": "bar",
				},
			}

			return dataExporter.Export(mockData)
		},
	}
}
