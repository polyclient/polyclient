// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package cli

import (
	"context"
	"os/exec"

	"github.com/polyclient/polyclient/internal/engine"
	"github.com/urfave/cli/v3"
)

// NewLogCommand creates a new command for viewing or managing logs.
func NewLogCommand(e *engine.Engine) *cli.Command {
	return &cli.Command{
		Name:  "log",
		Usage: "View or manage logs",
		Action: func(context.Context, *cli.Command) error {
			return exec.Command("less", "/tmp/polyclient.log").Start() // TODO: improve this
		},
	}
}
