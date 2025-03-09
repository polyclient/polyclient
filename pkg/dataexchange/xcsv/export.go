// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package xcsv

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/samber/lo"
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

// WithDateFormat sets a custom date format.
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

// Format writes the provided data to the given writer in CSV format.
//
// It supports two formats for `data`:
//  1. A flat slice of strings (`[]string`), treated as a single-column CSV.
//  2. A slice of maps (`[]map[string]any`), where keys become headers and values form rows.
//
// If the input is empty, nothing is written. If the data format is invalid, an error is returned.
func (ex *CsvExporter) Export(w io.Writer, data any) error {
	switch v := data.(type) {
	case []any:
		if len(v) == 0 {
			return nil
		}

		return ex.formatSlice(w, v)
	default:
		return fmt.Errorf("unsupported data type: %T", data)
	}
}

// formatSlice formats and writes the provided slice to the writer in CSV format.
func (ex *CsvExporter) formatSlice(w io.Writer, data []any) error {
	writer := csv.NewWriter(w)
	writer.Comma = ex.Comma

	if ex.UseCRLF {
		writer.UseCRLF = true
	}

	defer writer.Flush()

	switch data[0].(type) {
	case map[string]any:
		return ex.formatMapSlice(writer, data)
	default:
		return ex.formatSingleColumnSlice(writer, data)
	}
}

// formatMapSlice writes `[]map[string]any` as a multi-column CSV.
func (ex *CsvExporter) formatMapSlice(w *csv.Writer, data []any) error {
	first, ok := data[0].(map[string]any)
	if !ok {
		return fmt.Errorf("first element is not a map: %T", data[0])
	}

	headers := lo.Uniq(lo.Keys(first))

	if err := w.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	for _, item := range data {
		record, ok := item.(map[string]any)
		if !ok {
			return fmt.Errorf("item is not a map: %T", item)
		}

		row := lo.Map(headers, func(header string, _ int) string {
			return ex.formatValue(record[header])
		})

		if err := w.Write(row); err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	return nil
}

// formatSingleColumnCSV writes `[]any` as a single-column CSV.
func (ex *CsvExporter) formatSingleColumnSlice(w *csv.Writer, data []any) error {
	rows := lo.Map(data, func(item any, _ int) []string {
		return []string{ex.formatValue(item)}
	})

	if err := w.WriteAll(rows); err != nil {
		return fmt.Errorf("failed to write rows: %w", err)
	}

	return nil
}

// formatValue converts any value to a string, respecting the specified date format.
func (ex *CsvExporter) formatValue(value any) string {
	switch v := value.(type) {
	case time.Time:
		return v.Format(ex.DateFormat)
	case *time.Time:
		if v == nil {
			return ""
		}

		return v.Format(ex.DateFormat)
	case fmt.Stringer:
		return v.String()
	case error:
		return v.Error()
	default:
		return fmt.Sprintf("%v", value)
	}
}
