// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package xcsv_test

import (
	"bytes"
	"encoding/csv"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/polyclient/polyclient/pkg/dataexchange/xcsv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCsvExporter_Defaults(t *testing.T) {
	t.Parallel()

	ex := xcsv.NewCsvExporter()
	assert.Equal(t, ',', ex.Comma)
	assert.False(t, ex.UseCRLF)
	assert.Equal(t, time.RFC3339, ex.DateFormat)
}

func TestNewCsvExporter_WithOptions(t *testing.T) {
	t.Parallel()

	ex := xcsv.NewCsvExporter(
		xcsv.WithComma(';'),
		xcsv.WithCRLF(true),
		xcsv.WithDateFormat("2006-01-02"),
	)

	assert.Equal(t, ';', ex.Comma)
	assert.True(t, ex.UseCRLF)
	assert.Equal(t, "2006-01-02", ex.DateFormat)
}

func TestCsvExporter_Export_SingleColumn(t *testing.T) {
	t.Parallel()

	data := []any{"uno", "dos", "tres"}

	var buf bytes.Buffer

	ex := xcsv.NewCsvExporter()
	err := ex.Export(&buf, data)
	require.NoError(t, err)

	r := csv.NewReader(&buf)
	records, err := r.ReadAll()
	require.NoError(t, err)

	expected := [][]string{{"uno"}, {"dos"}, {"tres"}}
	assert.Equal(t, expected, records)
}

func TestCsvExporter_Export_MapSlice(t *testing.T) {
	t.Parallel()

	data := []any{
		map[string]any{"name": "Juan", "age": 30},
		map[string]any{"name": "Maria", "age": 25},
	}

	var buf bytes.Buffer

	ex := xcsv.NewCsvExporter()
	err := ex.Export(&buf, data)
	require.NoError(t, err)

	r := csv.NewReader(&buf)
	records, err := r.ReadAll()
	require.NoError(t, err)

	expected := [][]string{
		{"name", "age"},
		{"Juan", "30"},
		{"Maria", "25"},
	}
	assert.Equal(t, expected, records)
}

func TestCsvExporter_Export_EmptyInput(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	ex := xcsv.NewCsvExporter()
	err := ex.Export(&buf, []any{})
	assert.NoError(t, err)
	assert.Empty(t, buf.String())
}

func TestCsvExporter_Export_UnsupportedType(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	ex := xcsv.NewCsvExporter()
	err := ex.Export(&buf, 123)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported data type")
}

func TestCsvExporter_Export_InvalidMapSlice(t *testing.T) {
	t.Parallel()

	data := []any{
		map[string]any{"name": "Alice", "age": 30},
		"invalid", // Should cause an error
	}

	var buf bytes.Buffer

	ex := xcsv.NewCsvExporter()
	err := ex.Export(&buf, data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "item is not a map")
}

func TestCsvExporter_FormatValue(t *testing.T) {
	t.Parallel()

	ex := xcsv.NewCsvExporter(xcsv.WithDateFormat("2006-01-02"))
	timeVal := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{"string", "hello", "hello"},
		{"integer", 42, "42"},
		{"float", 3.14, "3.14"},
		{"bool", true, "true"},
		{"time", timeVal, "2023-01-01"},
		{"time pointer", &timeVal, "2023-01-01"},
		{"error", errors.New("test error"), "test error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			err := ex.Export(&buf, []any{tt.input})
			require.NoError(t, err)

			lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
			assert.Equal(t, 1, len(lines), "should have exactly one output row")
			assert.Equal(t, tt.expected, lines[0], "formatted value should match expected output")
		})
	}
}
