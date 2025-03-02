// Copyright (C) 2025 Juan Mesa and contributors
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3
// as published by the Free Software Foundation, with the Runtime
// Library Exception. See the COPYING.RUNTIME file for details.

package commands

import (
	"context"
	"fmt"

	"github.com/polyclient/polyclient/internal/version"
	"github.com/urfave/cli/v3"
)

// NewVersionCommand creates a CLI command for displaying detailed version information.
func NewVersionCommand() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Show more detailed version information",
		Action: func(context.Context, *cli.Command) error {
			//nolint:forbidigo  // Printing version info to stdout is intentional here.
			fmt.Println(version.String())

			return nil
		},
	}
}
