// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package gui

import (
	"embed"
	_ "embed"
	"io/fs"
)

//go:embed all:dist
var distDir embed.FS

// DistDirFS is the embedded dist directory as a file system for use with http.FS
var DistDirFS, _ = fs.Sub(distDir, "dist")
