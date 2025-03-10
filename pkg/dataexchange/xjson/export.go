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

// Export writes a slice to JSON, supporting primitive types, structs, and maps.
func (ex *JsonExporter) Export(w io.Writer, data any) error {
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
	enc.SetIndent("", ex.IndentString)
	enc.SetEscapeHTML(ex.EscapeHTML)

	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
