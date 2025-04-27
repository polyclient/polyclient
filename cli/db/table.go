package db

import (
	"context"
	"fmt"

	"github.com/polyclient/polyclient/internal/db"
	"github.com/polyclient/polyclient/internal/engine"
	"github.com/urfave/cli/v3"
)

func newTableCommand(e *engine.Engine) *cli.Command {
	return &cli.Command{
		Name:  "table",
		Usage: "Manage database tables",
		Commands: []*cli.Command{
			newTableListCommand(e),
			newTableGetCommand(e),
		},
	}
}

func newTableListCommand(e *engine.Engine) *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all database tables",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "schema",
				Aliases: []string{"s"},
				Usage:   "Filter tables by schema",
			},
			&cli.StringFlag{
				Name:    "filter",
				Aliases: []string{"f"},
				Usage:   "Filter tables by name",
			},
			&cli.IntFlag{
				Name:    "limit",
				Aliases: []string{"n"},
				Usage:   "Maximum number of tables to return",
				Value:   100,
			},
			&cli.IntFlag{
				Name:    "offset",
				Aliases: []string{"skip"},
				Usage:   "Offset for pagination",
				Value:   0,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			connName, err := getConnectionName(ctx, cmd)
			if err != nil {
				return fmt.Errorf("failed to get connection name: %w", err)
			}

			flagSchema := cmd.String("schema")
			flagFilter := cmd.String("filter")
			flagLimit := cmd.Int("limit")
			flagOffset := cmd.Int("offset")

			tables, err := e.SDK.Inspector().ListTables(ctx, connName,
				db.WithTablesSchema(flagSchema),
				db.WithTablesFilter(flagFilter),
				db.WithTablesLimit(int(flagLimit)),
				db.WithTablesOffset(int(flagOffset)),
			)
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

func newTableGetCommand(e *engine.Engine) *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "Get details of a table in a database",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "schema",
				Aliases: []string{"s"},
				Usage:   "Table schema",
			},
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Table name",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			connName, err := getConnectionName(ctx, cmd)
			if err != nil {
				return fmt.Errorf("failed to get connection name: %w", err)
			}

			flagSchema := cmd.String("schema")
			flagName := cmd.String("name")

			table, err := e.SDK.Inspector().GetTable(ctx, connName, flagName,
				db.WithTableSchema(flagSchema),
			)
			if err != nil {
				return fmt.Errorf("failed to get table: %w", err)
			}

			fmt.Printf("Name: %s\n", table.Name)
			fmt.Printf("Schema: %s\n", table.Schema)
			fmt.Printf("Columns:\n")

			for i := range table.Columns {
				fmt.Printf("  %s\n", table.Columns[i].Name)
			}

			return nil
		},
	}
}
