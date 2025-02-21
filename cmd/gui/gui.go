// Copyright (C) 2025 Juan Mesa and contributors
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3
// as published by the Free Software Foundation, with the Runtime
// Library Exception. See the COPYING.RUNTIME file for details.

package gui

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func NewGuiCommand() *cli.Command {
	return &cli.Command{
		Name:     "gui",
		Usage:    "Launch PolyClient in GUI mode",
		Category: "GUI",
		Action: func(context.Context, *cli.Command) error {
			fmt.Println("Launching PolyClient GUI")
			return nil
		},
	}
}
