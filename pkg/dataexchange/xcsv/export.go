// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package xcsv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/polyclient/polyclient/pkg/stringify"
)

// CsvExporter is a data exporter for CSV format.
type CsvExporter struct {
	// Comma is the character used to separate fields (default ',').
	Comma rune
	// UseCRLF indicates whether to use CRLF line endings (default false).
	UseCRLF bool
	// DateFormat is the format for date fields (default time.RFC3339).
	DateFormat string
}

// CsvExporterOption defines a functional option for configuring CSV output.
type CsvExporterOption func(*CsvExporter)

// WithComma sets a custom field separator.
func WithComma(comma rune) CsvExporterOption {
	return func(opts *CsvExporter) {
		opts.Comma = comma
	}
}

// WithCRLF sets whether to use CRLF line endings.
func WithCRLF(useCRLF bool) CsvExporterOption {
	return func(opts *CsvExporter) {
		opts.UseCRLF = useCRLF
	}
}

// WithDateFormat sets a custom date format for date fields.
func WithDateFormat(format string) CsvExporterOption {
	return func(opts *CsvExporter) {
		opts.DateFormat = format
	}
}

// NewCsvExporter creates a new CsvExporter with the specified options.
func NewCsvExporter(opts ...CsvExporterOption) *CsvExporter {
	ex := &CsvExporter{
		Comma:      ',',
		UseCRLF:    false,
		DateFormat: time.RFC3339,
	}

	for _, opt := range opts {
		opt(ex)
	}

	return ex
}

// Export writes the provided slice data to the given writer in CSV format.
// The data must be a slice of one of:
// - primitive types (writes single column)
// - structs (writes headers from exported fields)
// - maps[string]any (writes headers from first row's keys)
// Returns an error if data is not a slice or if writing fails.
func (ex *CsvExporter) Export(w io.Writer, data any) error {
	if w == nil {
		return errors.New("writer cannot be nil")
	}

	v := reflect.ValueOf(data)

	if v.Kind() != reflect.Slice {
		return fmt.Errorf("expected a slice, got %T", data)
	}

	if v.Len() == 0 {
		return nil
	}

	converted := make([]any, v.Len())
	for i := range v.Len() {
		converted[i] = v.Index(i).Interface()
	}

	return ex.formatSlice(w, converted)
}

// formatSlice formats the data as CSV, detecting the appropriate format based on the first element's type.
func (ex *CsvExporter) formatSlice(w io.Writer, data []any) error {
	writer := csv.NewWriter(w)
	writer.Comma = ex.Comma
	writer.UseCRLF = ex.UseCRLF

	defer writer.Flush()

	switch first := data[0].(type) {
	case map[string]any:
		return ex.formatMapSlice(writer, data)
	default:
		if reflect.TypeOf(first).Kind() == reflect.Struct {
			return ex.formatStructSlice(writer, data)
		}

		return ex.formatSingleColumnSlice(writer, data)
	}
}

// formatMapSlice writes `[]map[string]any` as a multi-column CSV.
func (ex *CsvExporter) formatMapSlice(w *csv.Writer, data []any) error {
	first, ok := data[0].(map[string]any)
	if !ok {
		return fmt.Errorf("first element is not a map: %T", data[0])
	}

	var headers = make([]string, 0, len(first))
	for header := range first {
		headers = append(headers, header)
	}

	if err := w.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for _, item := range data {
		record, ok := item.(map[string]any)
		if !ok {
			return fmt.Errorf("item is not a map: %T", item)
		}

		row := make([]string, 0, len(headers))

		for _, header := range headers {
			row = append(row, stringify.Stringify(record[header]))
		}

		if err := w.Write(row); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

// formatStructSlice writes `[]struct` as a multi-column CSV.
func (ex *CsvExporter) formatStructSlice(w *csv.Writer, data []any) error {
	var headers []string

	first := reflect.TypeOf(data[0])
	for i := range first.NumField() {
		field := first.Field(i)
		if field.PkgPath == "" {
			headers = append(headers, field.Name)
		}
	}

	if err := w.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for _, item := range data {
		record := reflect.ValueOf(item)

		var row []string

		for _, header := range headers {
			field := record.FieldByName(header)
			if field.IsValid() && field.CanInterface() {
				row = append(row, stringify.Stringify(field.Interface()))
			} else {
				row = append(row, "")
			}
		}

		if err := w.Write(row); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

// formatSingleColumnSlice writes `[]any` as a single-column CSV.
func (ex *CsvExporter) formatSingleColumnSlice(w *csv.Writer, data []any) error {
	rows := make([][]string, 0, len(data))

	for _, item := range data {
		rows = append(rows, []string{stringify.Stringify(item)})
	}

	if err := w.WriteAll(rows); err != nil {
		return fmt.Errorf("failed to write rows: %w", err)
	}

	return nil
}
