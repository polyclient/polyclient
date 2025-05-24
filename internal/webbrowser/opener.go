// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package webbrowser

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
)

var defaultRunCommand = func(ctx context.Context, name string, args ...string) error {
	return exec.CommandContext(ctx, name, args...).Run()
}

var defaultGOOS = func() string {
	return runtime.GOOS
}

var defaultIsWSL = func() bool {
	return false
}

// OpenURL opens the specified URL in the default web browser.
func OpenURL(ctx context.Context, rawURL string) error {
	if rawURL == "" {
		return errors.New("URL cannot be empty")
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	os := defaultGOOS()

	if defaultIsWSL() {
		return defaultRunCommand(ctx, "powershell.exe", "Start-Process", parsedURL.String())
	}

	switch os {
	case "linux", "freebsd", "openbsd":
		return defaultRunCommand(ctx, "xdg-open", parsedURL.String())
	case "darwin":
		return defaultRunCommand(ctx, "open", parsedURL.String())
	case "windows":
		return defaultRunCommand(ctx, "cmd.exe", "/c", "start", "", parsedURL.String())
	default:
		return fmt.Errorf("unsupported operating system: %s - try manually opening a browser", os)
	}
}
