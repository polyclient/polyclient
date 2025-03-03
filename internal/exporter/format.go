// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package exporter

import (
	"io"
)

// Format represents supported export formats.
type Format string

const (
	CSV      = "csv"
	TSV      = "tsv"
	JSON     = "json"
	HTML     = "html"
	XML      = "xml"
	YAML     = "yaml"
	TOML     = "toml"
	Markdown = "markdown"
)

// FormatDetails contains metadata and handlers for each format.
type FormatDetails struct {
	MIMEType  string
	FileExt   string
	Formatter func(io.Writer, any) error
}

var formatRegistry = map[Format]FormatDetails{
	CSV: {
		MIMEType:  "text/csv",
		FileExt:   "csv",
		Formatter: FormatCSV,
	},
	TSV: {
		MIMEType:  "text/tab-separated-values",
		FileExt:   "tsv",
		Formatter: FormatTSV,
	},
	JSON: {
		MIMEType:  "application/json",
		FileExt:   "json",
		Formatter: FormatJSON,
	},
	HTML: {
		MIMEType:  "text/html",
		FileExt:   "html",
		Formatter: FormatHTML,
	},
	XML: {
		MIMEType:  "application/xml",
		FileExt:   "xml",
		Formatter: FormatXML,
	},
	YAML: {
		MIMEType:  "application/yaml",
		FileExt:   "yaml",
		Formatter: FormatYAML,
	},
	TOML: {
		MIMEType:  "application/toml",
		FileExt:   "toml",
		Formatter: FormatTOML,
	},
	Markdown: {
		MIMEType:  "text/markdown",
		FileExt:   "md",
		Formatter: FormatMarkdown,
	},
}

// ValidateFormat returns true if the format is supported, otherwise false.
func ValidateFormat(format Format) bool {
	_, ok := formatRegistry[format]
	return ok
}

// GetSupportedFormats returns a slice of supported formats.
func GetSupportedFormats() []string {
	formats := make([]string, 0, len(formatRegistry))

	for format := range formatRegistry {
		formats = append(formats, string(format))
	}

	return formats
}
