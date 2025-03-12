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

// JSONExporter is a data exporter for JSON format.
type JSONExporter struct {
	// IndentString is the indentation string used in the JSON output (default: "  ").
	IndentString string
	// EscapeHTML specifies whether to escape HTML characters in the JSON output
	// (default: true -> converts &, <, and > to \u0026, \u003c, and \u003e).
	EscapeHTML bool
}

// JSONExporterOption is a functional option for configuring JsonExporter.
type JSONExporterOption func(*JSONExporter)

// WithIndentString sets the indentation string for JsonExporter.
func WithIndentString(indent string) JSONExporterOption {
	return func(exp *JSONExporter) {
		exp.IndentString = indent
	}
}

// WithEscapeHTML sets whether to escape HTML characters in JsonExporter.
func WithEscapeHTML(escape bool) JSONExporterOption {
	return func(exp *JSONExporter) {
		exp.EscapeHTML = escape
	}
}

// NewJSONExporter creates a new instance of JsonExporter.
func NewJSONExporter(opts ...JSONExporterOption) *JSONExporter {
	exporter := &JSONExporter{
		IndentString: "  ",
		EscapeHTML:   true,
	}

	for _, opt := range opts {
		opt(exporter)
	}

	return exporter
}

// Export writes a slice to JSON, supporting primitive types, structs, and maps.
func (exp *JSONExporter) Export(w io.Writer, data any) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}

	v := reflect.ValueOf(data)

	if !v.IsValid() || (v.Kind() == reflect.Pointer && v.IsNil()) {
		return errors.New("data cannot be nil")
	}

	// Handle empty slices explicitly to ensure consistent output
	if v.Kind() == reflect.Slice && v.Len() == 0 {
		_, err := w.Write([]byte("[]\n"))
		return err
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", exp.IndentString)
	enc.SetEscapeHTML(exp.EscapeHTML)

	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
