// Copyright (C) 2025 Juan Mesa and contributors
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3
// as published by the Free Software Foundation, with the Runtime
// Library Exception. See the COPYING.RUNTIME file for details.

package version

import "fmt"

// build holds the build metadata and is populated at build time.
var (
	version     = "dev"     // Git tag (e.g., "v1.2.3")
	branch      = "unknown" // Git branch (e.g., "master")
	shortCommit = "none"    // Git short commit SHA (e.g., "0fd6153")
	fullCommit  = "none"    // Git full commit SHA (e.g., "0fd6153429327455ec5bca2cda839116cfcb6a19")
	date        = "unknown" // Build date in RFC3339 format (e.g., "2025-03-02T18:46:18Z")
	os          = "unknown" // Operating system (e.g., "linux")
	arch        = "unknown" // CPU architecture (e.g., "amd64")
)

// String returns a formatted multi-line string that displays the build metadata, including version, branch, short and full commit SHAs, build date, operating system, and CPU architecture.
func String() string {
	return fmt.Sprintf(
		"Version: %s\nBranch: %s\nShort commit: %s\nFull commit: %s\nBuild date: %s\nOS: %s\nArch: %s",
		version, branch, shortCommit, fullCommit, date, os, arch)
}

// Version returns the build version of the application. This value is typically set at compile time.
func Version() string {
	return version
}

// Branch returns the Git branch name used during the build.
func Branch() string {
	return branch
}

// ShortCommit returns the abbreviated SHA of the commit used to build the software.
func ShortCommit() string {
	return shortCommit
}

// FullCommit returns the full Git commit SHA from which the binary was built.
func FullCommit() string {
	return fullCommit
}

// Date returns the build date in RFC3339 format.
func Date() string {
	return date
}

// OS returns the operating system build metadata.
// It retrieves the operating system name set during build time, defaulting to "unknown" if not specified.
func OS() string {
	return os
}

// Arch returns the CPU architecture recorded at build time.
func Arch() string {
	return arch
}
