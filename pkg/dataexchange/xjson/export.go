// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package xjson

import (
	"encoding/json"
	"io"
)

// JsonExporter is a data exporter for JSON format.
type JsonExporter struct {
	// IndentString is the indentation string used in the JSON output (default: "  ").
	IndentString string
	// EscapeHTML specifies whether to escape HTML characters in the JSON output (default: true -> converts &, <, and > to \u0026, \u003c, and \u003e).
	EscapeHTML bool
}

// JsonExporterOption is a functional option for configuring JsonExporter.
type JsonExporterOption func(*JsonExporter)

// WithIndentString sets the indentation string for JsonExporter.
func WithIndentString(indent string) JsonExporterOption {
	return func(ex *JsonExporter) {
		ex.IndentString = indent
	}
}

// WithEscapeHTML sets whether to escape HTML characters in JsonExporter.
func WithEscapeHTML(escape bool) JsonExporterOption {
	return func(ex *JsonExporter) {
		ex.EscapeHTML = escape
	}
}

// NewJsonExporter creates a new instance of JsonExporter.
func NewJsonExporter(opts ...JsonExporterOption) *JsonExporter {
	ex := &JsonExporter{
		IndentString: "  ",
		EscapeHTML:   true,
	}

	for _, opt := range opts {
		opt(ex)
	}

	return ex
}

// FormatJSON formats and writes the provided data to the writer in JSON format.
func (ex *JsonExporter) Export(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", ex.IndentString)
	enc.SetEscapeHTML(ex.EscapeHTML)

	return enc.Encode(data)
}
