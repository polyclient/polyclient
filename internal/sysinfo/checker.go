// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package sysinfo

import (
	"os"
	"path"
	"runtime"
	"slices"
	"strings"
)

var defaultFileReader = func(name string) ([]byte, error) {
	return os.ReadFile(path.Clean(name))
}

var defaultGOOS = func() string {
	return runtime.GOOS
}

// IsWSL checks if the system is running in Windows Subsystem for Linux (WSL).
func IsWSL() bool {
	if defaultGOOS() != "linux" {
		return false
	}

	lookupPaths := []string{
		path.Join("/proc", "version"),
		path.Join("/proc", "sys", "kernel", "osrelease"),
	}

	return slices.ContainsFunc(lookupPaths, fileContainsWSLIndicator)
}

// fileContainsWSLIndicators checks if a file contain WSL indicators.
func fileContainsWSLIndicator(filePath string) bool {
	data, err := defaultFileReader(filePath)
	if err != nil {
		return false
	}

	content := strings.ToLower(string(data))

	return strings.Contains(content, "microsoft") || strings.Contains(content, "wsl")
}
