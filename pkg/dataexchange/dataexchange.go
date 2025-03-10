// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package dataexchange

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"slices"

	"github.com/polyclient/polyclient/pkg/dataexchange/xcsv"
	"github.com/polyclient/polyclient/pkg/dataexchange/xhtml"
	"github.com/polyclient/polyclient/pkg/dataexchange/xjson"
)

// Format represents supported export formats.
type Format string

const (
	FormatCsv      = "csv"
	FormatTsv      = "tsv"
	FormatJson     = "json"
	FormatHtml     = "html"
	FormatXml      = "xml"
	FormatYaml     = "yaml"
	FormatToml     = "toml"
	FormatMarkdown = "markdown"
)

// Importer represents the import data format interface.
type Importer interface {
	Import(r io.Reader) (any, error)
}

// Exporter represents the export data format interface.
type Exporter interface {
	Export(w io.Writer, data any) error
}

// ExporterEntry contains details for a specific export format.
type Entry struct {
	MIMEType string
	FileExt  string
	Importer
	Exporter
}

// registry maps supported import and export formats to their details.
var registry = map[Format]Entry{
	FormatCsv: {
		MIMEType: "text/csv",
		FileExt:  "csv",
		Exporter: xcsv.NewCsvExporter(),
	},
	FormatTsv: {
		MIMEType: "text/tab-separated-values",
		FileExt:  "tsv",
		Exporter: xcsv.NewCsvExporter(xcsv.WithComma('\t')),
	},
	FormatJson: {
		MIMEType: "application/json",
		FileExt:  "json",
		Exporter: xjson.NewJsonExporter(),
	},
	FormatHtml: {
		MIMEType: "text/html",
		FileExt:  "html",
		Exporter: xhtml.NewHtmlExporter(),
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

// GetSupportedFormats returns a list of supported export formats.
func GetSupportedExportFormats() []Format {
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
