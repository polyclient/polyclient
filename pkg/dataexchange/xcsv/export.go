// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

// Package xcsv exports data as CSV.
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

// CSVExporter exports data as CSV.
type CSVExporter struct {
	Comma      rune
	UseCRLF    bool
	DateFormat string
}

// CSVExporterOption is a functional option type for configuring a CsvExporter.
type CSVExporterOption func(*CSVExporter)

func WithComma(comma rune) CSVExporterOption {
	return func(opts *CSVExporter) { opts.Comma = comma }
}

func WithCRLF(useCRLF bool) CSVExporterOption {
	return func(opts *CSVExporter) { opts.UseCRLF = useCRLF }
}

func WithDateFormat(format string) CSVExporterOption {
	return func(opts *CSVExporter) { opts.DateFormat = format }
}

func NewCSVExporter(opts ...CSVExporterOption) *CSVExporter {
	ex := &CSVExporter{
		Comma:      ',',
		UseCRLF:    false,
		DateFormat: time.RFC3339,
	}

	for _, opt := range opts {
		opt(ex)
	}

	return ex
}

// Export writes a slice to CSV, supporting primitive types, structs, and maps.
func (ex *CSVExporter) Export(w io.Writer, data any) error {
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
func (ex *CSVExporter) formatSlice(w io.Writer, v reflect.Value) error {
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
func (*CSVExporter) formatMapSlice(w *csv.Writer, v reflect.Value) error {
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
func (*CSVExporter) formatStructSlice(w *csv.Writer, v reflect.Value) error {
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
func (*CSVExporter) formatSingleColumnSlice(w *csv.Writer, v reflect.Value) error {
	rows := make([][]string, v.Len())
	for i := range v.Len() {
		rows[i] = []string{stringify.Stringify(v.Index(i).Interface())}
	}

	if err := w.WriteAll(rows); err != nil {
		return fmt.Errorf("failed to write rows: %w", err)
	}

	return nil
}
