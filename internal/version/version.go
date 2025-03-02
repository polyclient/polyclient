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

func String() string {
	return fmt.Sprintf(
		"Version: %s\nBranch: %s\nShort commit: %s\nFull commit: %s\nBuild date: %s\nOS: %s\nArch: %s",
		version, branch, shortCommit, fullCommit, date, os, arch)
}

func Version() string {
	return version
}

func Branch() string {
	return branch
}

func ShortCommit() string {
	return shortCommit
}

func FullCommit() string {
	return fullCommit
}

func Date() string {
	return date
}

func OS() string {
	return os
}

func Arch() string {
	return arch
}
