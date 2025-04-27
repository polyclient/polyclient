package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/polyclient/polyclient/internal/db"
	"github.com/polyclient/polyclient/internal/engine"
	"github.com/urfave/cli/v3"
)

func newConnectionCommand(e *engine.Engine) *cli.Command {
	return &cli.Command{
		Name:  "connection",
		Usage: "Manage database connections",
		Commands: []*cli.Command{
			newConnectionCreateCommand(e),
			newConnectionListCommand(e),
			newConnectionDeleteCommand(e),
			newConnectionPingCommand(e),
		},
	}
}

func newConnectionCreateCommand(e *engine.Engine) *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create a new database connection",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "driver",
				Aliases:  []string{"d"},
				Usage:    "Driver of the connection",
				Required: true,
				Config:   cli.StringConfig{TrimSpace: true},
			},
			&cli.StringFlag{
				Name:     "name",
				Aliases:  []string{"n"},
				Usage:    "Name of the connection",
				Required: true,
				Config:   cli.StringConfig{TrimSpace: true},
			},
			&cli.StringFlag{
				Name:     "dsn",
				Usage:    "Data source name of the connection",
				Required: true,
				Config:   cli.StringConfig{TrimSpace: true},
			},
			&cli.StringFlag{
				Name:   "color-tag",
				Usage:  "Color tag of the connection",
				Config: cli.StringConfig{TrimSpace: true},
			},
			&cli.BoolFlag{
				Name:  "save-creds",
				Usage: "Save credentials for the connection",
			},
			&cli.BoolFlag{
				Name:  "pinned",
				Usage: "Pin the connection",
			},
			&cli.BoolFlag{
				Name:  "confirm-before-connect",
				Usage: "Confirm before connecting to the database",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			connName, _ := getConnectionName(ctx, cmd)
			if connName != "" {
				return errors.New("connection must not be specified when creating connections")
			}

			flagDriver := cmd.String("driver")
			flagName := cmd.String("name")
			flagColorTag := cmd.String("color-tag")
			flagSaveCreds := cmd.Bool("save-creds")
			flagPinned := cmd.Bool("pinned")
			flagConfirmBeforeConnect := cmd.Bool("confirm-before-connect")
			flagDSN := cmd.String("dsn")

			profile := db.ConnectionProfile{
				Driver:               flagDriver,
				Name:                 flagName,
				ColorTag:             flagColorTag,
				SaveCreds:            flagSaveCreds,
				Pinned:               flagPinned,
				ConfirmBeforeConnect: flagConfirmBeforeConnect,
				Config:               db.ConnectionConfig{"dsn": flagDSN},
				CreatedAt:            time.Now(),
				LastUsedAt:           time.Now(),
			}

			connectionStore := e.SDK.GetManager().GetStore()

			if err := connectionStore.SaveProfile(ctx, &profile); err != nil {
				return fmt.Errorf("failed to save profile: %w", err)
			}

			return nil
		},
	}
}

func newConnectionListCommand(e *engine.Engine) *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "List all database connections",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			connName, _ := getConnectionName(ctx, cmd)
			if connName != "" {
				return errors.New("connection must not be specified when listing connections")
			}

			profiles, err := e.SDK.GetManager().GetStore().ListProfiles(ctx)
			if err != nil {
				return fmt.Errorf("failed to list profiles: %w", err)
			}

			for _, profile := range profiles {
				fmt.Printf("%s (%s)\n", profile.Name, profile.Driver)
			}

			return nil
		},
	}
}

func newConnectionDeleteCommand(e *engine.Engine) *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "Delete a database connection",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			connName, err := getConnectionName(ctx, cmd)
			if err != nil {
				return fmt.Errorf("failed to get connection name: %w", err)
			}

			if err := e.SDK.GetManager().GetStore().DeleteProfile(ctx, connName); err != nil {
				return fmt.Errorf("failed to delete profile: %w", err)
			}

			return nil
		},
	}
}

func newConnectionPingCommand(e *engine.Engine) *cli.Command {
	return &cli.Command{
		Name:  "ping",
		Usage: "Ping a database to check connectivity",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			connName, err := getConnectionName(ctx, cmd)
			if err != nil {
				return fmt.Errorf("failed to get connection name: %w", err)
			}

			err = e.SDK.Ping(ctx, connName)
			if err != nil {
				return fmt.Errorf("failed to open connection: %w", err)
			}

			infoDB, err := e.SDK.Info().CurrentDatabase(ctx, connName)
			if err != nil {
				return fmt.Errorf("failed to get current database: %w", err)
			}

			infoVersion, err := e.SDK.Info().ServerVersion(ctx, connName)
			if err != nil {
				return fmt.Errorf("failed to get server version: %w", err)
			}

			fmt.Printf("Database: %s\n", infoDB)
			fmt.Printf("Server Version: %s\n", infoVersion)

			return nil
		},
	}
}
