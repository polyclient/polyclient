// Copyright (C) 2025 Juan Mesa and contributors
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3
// as published by the Free Software Foundation, with the Runtime
// Library Exception. See the COPYING.RUNTIME file for details.

package db

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/polyclient/polyclient/internal/exporter"
	"github.com/urfave/cli/v3"
)

// NewQueryCommand returns a new query command that can be used to query a database from the CLI.
func NewQueryCommand() *cli.Command {
	return &cli.Command{
		Name:     "execute",
		Usage:    "Execute a query against a database (SQL or NoSQL)",
		Category: "Database",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "command",
				Aliases:  []string{"c"},
				Usage:    "Command to execute (e.g., 'SELECT * FROM users')",
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
			// command := cmd.String("command")
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
				Output: w,
			})

			mockData := []struct {
				Key   string
				Value any
			}{
				{
					Key: "User1", Value: struct {
						name string
						age  int
					}{
						name: "John",
						age:  30,
					},
				},
				{
					Key: "User2", Value: struct {
						name string
						age  int
					}{
						name: "Jane",
						age:  25,
					},
				},
			}

			return dataExporter.Export(mockData)
		},
	}
}
