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
	"sort"
	"time"

	"github.com/polyclient/polyclient/pkg/stringify"
)

// CsvExporter exports data as CSV.
type CsvExporter struct {
	Comma      rune
	UseCRLF    bool
	DateFormat string
}

type CsvExporterOption func(*CsvExporter)

func WithComma(comma rune) CsvExporterOption {
	return func(opts *CsvExporter) { opts.Comma = comma }
}

func WithCRLF(useCRLF bool) CsvExporterOption {
	return func(opts *CsvExporter) { opts.UseCRLF = useCRLF }
}

func WithDateFormat(format string) CsvExporterOption {
	return func(opts *CsvExporter) { opts.DateFormat = format }
}

func NewCsvExporter(opts ...CsvExporterOption) *CsvExporter {
	ex := &CsvExporter{Comma: ',', UseCRLF: false, DateFormat: time.RFC3339}
	for _, opt := range opts {
		opt(ex)
	}

	return ex
}

// Export writes a slice to CSV, supporting primitive types, structs, and maps.
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

	return ex.formatSlice(w, v)
}

// formatSlice processes a slice and determines the CSV format.
func (ex *CsvExporter) formatSlice(w io.Writer, v reflect.Value) error {
	writer := csv.NewWriter(w)
	writer.Comma = ex.Comma
	writer.UseCRLF = ex.UseCRLF

	defer writer.Flush()

	first := v.Index(0).Interface()

	switch first.(type) {
	case map[string]any:
		return ex.formatMapSlice(writer, v)
	default:
		if reflect.TypeOf(first).Kind() == reflect.Struct {
			return ex.formatStructSlice(writer, v)
		}

		return ex.formatSingleColumnSlice(writer, v)
	}
}

// formatMapSlice writes `[]map[string]any` as a CSV.
func (ex *CsvExporter) formatMapSlice(w *csv.Writer, v reflect.Value) error {
	first := v.Index(0).Interface().(map[string]any)
	headers := make([]string, 0, len(first))

	for k := range first {
		headers = append(headers, k)
	}

	sort.Strings(headers) // Ensure consistent column ordering

	if err := w.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for i := range v.Len() {
		record := v.Index(i).Interface().(map[string]any)
		row := make([]string, len(headers))

		for j, header := range headers {
			row[j] = stringify.Stringify(record[header])
		}

		if err := w.Write(row); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

// formatStructSlice writes `[]struct` as a CSV.
func (ex *CsvExporter) formatStructSlice(w *csv.Writer, v reflect.Value) error {
	t := v.Index(0).Type()

	var headers []string

	var fieldIndices []int

	// Filter out private fields
	for i := range t.NumField() {
		field := t.Field(i)
		if field.PkgPath == "" {
			headers = append(headers, field.Name)
			fieldIndices = append(fieldIndices, i)
		}
	}

	if err := w.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for i := range v.Len() {
		record := v.Index(i)
		row := make([]string, len(headers))

		for j, idx := range fieldIndices {
			field := record.Field(idx)
			if field.CanInterface() {
				row[j] = stringify.Stringify(field.Interface())
			}
		}

		if err := w.Write(row); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

// formatSingleColumnSlice writes `[]any` as a single-column CSV.
func (ex *CsvExporter) formatSingleColumnSlice(w *csv.Writer, v reflect.Value) error {
	rows := make([][]string, v.Len())
	for i := range v.Len() {
		rows[i] = []string{stringify.Stringify(v.Index(i).Interface())}
	}

	if err := w.WriteAll(rows); err != nil {
		return fmt.Errorf("failed to write rows: %w", err)
	}

	return nil
}
