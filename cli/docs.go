// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package cli

import (
	"context"
	"os/exec"
	"runtime"
	"strings"

	"github.com/urfave/cli/v3"
)

// NewDocsCommand creates a CLI command for opening the documentation website.
func NewDocsCommand() *cli.Command {
	return &cli.Command{
		Name:  "docs",
		Usage: "Open documentation website",
		Action: func(context.Context, *cli.Command) error {
			url := "https://polyclient.pages.dev"

			if runtime.GOOS == "windows" {
				return exec.Command("cmd", "/c", "start", url).Start()
			}

			if runtime.GOOS == "darwin" {
				return exec.Command("open", url).Start()
			}

			if isWSL() {
				return exec.Command("cmd", "/c", "start", url).Start()
			}

			return exec.Command("xdg-open", url).Start()
		},
	}
}

// isWSL checks if the Go program is running inside Windows Subsystem for Linux.
func isWSL() bool {
	releaseData, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}

	return strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
}
