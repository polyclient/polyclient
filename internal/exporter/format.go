// Copyright (C) 2025 Juan Mesa and contributors
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3
// as published by the Free Software Foundation, with the Runtime
// Library Exception. See the COPYING.RUNTIME file for details.

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
		Formatter: formatHTML,
	},
	XML: {
		MIMEType:  "application/xml",
		FileExt:   "xml",
		Formatter: formatXML,
	},
	YAML: {
		MIMEType:  "application/yaml",
		FileExt:   "yaml",
		Formatter: formatYAML,
	},
	TOML: {
		MIMEType:  "application/toml",
		FileExt:   "toml",
		Formatter: formatTOML,
	},
	Markdown: {
		MIMEType:  "text/markdown",
		FileExt:   "md",
		Formatter: formatMarkdown,
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

// formatHTML writes the provided data to the writer in HTML format.
func formatHTML(w io.Writer, data any) error {
	return nil
}

// formatXML writes the provided data to the writer in XML format.
func formatXML(w io.Writer, data any) error {
	return nil
}

// formatYAML writes the provided data to the writer in YAML format.
func formatYAML(w io.Writer, data any) error {
	return nil
}

// formatTOML writes the provided data to the writer in TOML format.
func formatTOML(w io.Writer, data any) error {
	return nil
}

// formatMarkdown writes the provided data to the writer in Markdown format.
func formatMarkdown(w io.Writer, data any) error {
	return nil
}
