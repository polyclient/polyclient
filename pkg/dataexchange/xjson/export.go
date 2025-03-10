// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package xjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
)

// JsonExporter is a data exporter for JSON format.
type JsonExporter struct {
	// IndentString is the indentation string used in the JSON output (default: "  ").
	IndentString string
	// EscapeHTML specifies whether to escape HTML characters in the JSON output
	// (default: true -> converts &, <, and > to \u0026, \u003c, and \u003e).
	EscapeHTML bool
}

// JsonExporterOption is a functional option for configuring JsonExporter.
type JsonExporterOption func(*JsonExporter)

// WithIndentString sets the indentation string for JsonExporter.
// Common values are "  " (two spaces) or "\t" (tab).
func WithIndentString(indent string) JsonExporterOption {
	return func(ex *JsonExporter) {
		ex.IndentString = indent
	}
}

// WithEscapeHTML sets whether to escape HTML characters in JsonExporter.
// When true, &, <, and > are escaped to \u0026, \u003c, and \u003e respectively.
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

// Export writes the provided data to the given writer in JSON format.
// The data can be any JSON-serializable type, including:
// - structs (exported fields only)
// - maps with string keys
// - slices and arrays
// - primitive types
// Returns an error if:
// - writer is nil
// - data cannot be marshaled to JSON
// - writing to the writer fails.
func (ex *JsonExporter) Export(w io.Writer, data any) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}

	if data == nil {
		return errors.New("data cannot be nil")
	}

	// Handle empty slices specially to ensure consistent output
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice && v.Len() == 0 {
		_, err := w.Write([]byte("[]\n"))
		return err
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", ex.IndentString)
	enc.SetEscapeHTML(ex.EscapeHTML)

	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
