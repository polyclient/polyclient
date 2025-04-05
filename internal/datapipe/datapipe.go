// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

// Package datapipe provide functions for exporting and importing data.
package datapipe

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"slices"

	"github.com/polyclient/polyclient/internal/datapipe/xcsv"
	"github.com/polyclient/polyclient/internal/datapipe/xhtml"
	"github.com/polyclient/polyclient/internal/datapipe/xjson"
)

// Format represents supported export formats.
type Format string

// String returns the string representation of the format.
func (f Format) String() string {
	return string(f)
}

const (
	// FormatCSV represents the CSV format.
	FormatCSV = "csv"
	// FormatTSV represents the TSV format.
	FormatTSV = "tsv"
	// FormatJSON represents the JSON format.
	FormatJSON = "json"
	// FormatHTML represents the HTML format.
	FormatHTML = "html"
	// FormatXML represents the XML format.
	FormatXML = "xml"
	// FormatYAML represents the YAML format.
	FormatYAML = "yaml"
	// FormatTOML represents the TOML format.
	FormatTOML = "toml"
	// FormatMarkdown represents the Markdown format.
	FormatMarkdown = "markdown"
)

// Importer represents the import data format interface.
type Importer interface {
	// Import reads data from a reader and returns a slice of type `[]any`.
	Import(r io.Reader) ([]any, error)
}

// Exporter represents the export data format interface.
type Exporter interface {
	// Export writes a slice to a writer.
	Export(w io.Writer, data any) error
}

// Entry contains details for a specific export format.
type Entry struct {
	MIMEType string
	FileExt  string
	Importer
	Exporter
}

// registry maps supported import and export formats to their details.
var registry = map[Format]Entry{
	FormatCSV: {
		MIMEType: "text/csv",
		FileExt:  "csv",
		Exporter: xcsv.NewCSVExporter(),
	},
	FormatTSV: {
		MIMEType: "text/tab-separated-values",
		FileExt:  "tsv",
		Exporter: xcsv.NewCSVExporter(xcsv.WithDelimiter('\t')),
	},
	FormatJSON: {
		MIMEType: "application/json",
		FileExt:  "json",
		Exporter: xjson.NewJSONExporter(),
	},
	FormatHTML: {
		MIMEType: "text/html",
		FileExt:  "html",
		Exporter: xhtml.NewHTMLExporter(),
	},
	// FormatXml: {
	// 	MIMEType:  "application/xml",
	// 	FileExt:   "xml",
	// 	Exporter: xxml.NewXmlExporter(),
	// },
	// FormatYaml: {
	// 	MIMEType:  "application/yaml",
	// 	FileExt:   "yaml",
	// 	Exporter:  xyaml.NewYamlExporter(),
	// },
	// FormatToml: {
	// 	MIMEType:  "application/toml",
	// 	FileExt:   "toml",
	// 	Exporter:  xtoml.NewTomlExporter(),
	// },
	// FormatMarkdown: {
	// 	MIMEType:  "text/markdown",
	// 	FileExt:   "md",
	// 	Exporter:  xmd.NewMarkdownExporter(),
	// },
}

// GetAvailableExportFormats returns a list of supported export formats as strings.
func GetAvailableExportFormats() []Format {
	return slices.Collect(maps.Keys(registry))
}

// GetRegistryEntry returns the details for a specific export
// format along with its associated exporter function.
func GetRegistryEntry(format Format) (Entry, bool) {
	entry, ok := registry[format]

	return entry, ok
}

// ParseDataFromBytes parses data bytes into a slice of type T.
// Requires a valid format with a registered exporter. Returns
// an error if the format is invalid. T must be a valid JSON type.
func ParseDataFromBytes[T any](data []byte, format Format) (T, error) {
	_, ok := registry[format]
	if !ok {
		var zero T

		return zero, fmt.Errorf("format %s is not supported", format)
	}

	var parsedData T
	if err := json.Unmarshal(data, &parsedData); err != nil {
		var zero T

		return zero, fmt.Errorf("failed to parse data: %w", err)
	}

	return parsedData, nil
}
