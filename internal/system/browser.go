// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package system

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
)

// OpenBrowser opens the specified URL in the default browser for the current platform.
func OpenBrowser(targetURL string) error {
	parsedRequestURL, err := url.ParseRequestURI(targetURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	parsedURL := parsedRequestURL.String()

	if runtime.GOOS == "windows" || IsWSL() {
		return exec.Command("cmd", "/c", "start", parsedURL).Start()
	}

	if runtime.GOOS == "darwin" {
		return exec.Command("open", parsedURL).Start()
	}

	if err := exec.Command("xdg-open", parsedURL).Start(); err != nil {
		return fmt.Errorf("failed to open browser with xdg-open: %w", err)
	}

	return nil
}

// IsWSL checks if the system is running under Windows Subsystem for Linux.
func IsWSL() bool {
	releaseData, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}

	return strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
}
