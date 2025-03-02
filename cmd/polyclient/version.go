// Copyright (C) 2025 Juan Mesa and contributors
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3
// as published by the Free Software Foundation, with the Runtime
// Library Exception. See the COPYING.RUNTIME file for details.

package main

import (
	"context"
	"fmt"

	"github.com/polyclient/polyclient/internal/version"
	"github.com/urfave/cli/v3"
)

func NewVersionCommand() *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Show more detailed version information",
		Action: func(context.Context, *cli.Command) error {
			//nolint:forbidigo
			fmt.Println(version.String())

			return nil
		},
	}
}
