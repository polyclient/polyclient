// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package cli

import (
	"context"
	"os/exec"

	"github.com/polyclient/polyclient/internal/application"
	"github.com/urfave/cli/v3"
)

func NewLogCommand(app *application.Application) *cli.Command {
	return &cli.Command{
		Name:  "log",
		Usage: "View or manage logs",
		Action: func(context.Context, *cli.Command) error {
			return exec.Command("less", "/tmp/polyclient.log").Start()
		},
	}
}
