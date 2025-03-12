// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package xcsv_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/polyclient/polyclient/internal/datapipe/xcsv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestPerson struct {
	Name     string
	Age      int
	Active   bool
	JoinedAt time.Time
}

type PrivateFieldStruct struct {
	Public  string
	private string
}

func TestNewCsvExporter(t *testing.T) {
	t.Parallel()

	t.Run("default configuration", func(t *testing.T) {
		t.Parallel()

		exporter := xcsv.NewCSVExporter()
		assert.NotNil(t, exporter)
	})

	t.Run("custom configuration", func(t *testing.T) {
		t.Parallel()

		exporter := xcsv.NewCSVExporter(
			xcsv.WithDelimiter(';'),
			xcsv.WithCRLF(true),
			xcsv.WithDateFormat("2006-01-02"),
		)
		assert.NotNil(t, exporter)
	})
}

func TestExport(t *testing.T) {
	t.Parallel()

	t.Run("nil writer", func(t *testing.T) {
		t.Parallel()

		exporter := xcsv.NewCSVExporter()
		err := exporter.Export(nil, []string{"data"})
		require.Error(t, err)

		assert.Contains(t, err.Error(), "writer")
	})

	t.Run("invalid input - non-slice", func(t *testing.T) {
		t.Parallel()

		exporter := xcsv.NewCSVExporter()

		var buf bytes.Buffer

		err := exporter.Export(&buf, "not a slice")
		require.Error(t, err)
	})

	t.Run("empty slice", func(t *testing.T) {
		t.Parallel()

		exporter := xcsv.NewCSVExporter()

		var buf bytes.Buffer

		err := exporter.Export(&buf, []string{})
		require.NoError(t, err)

		assert.Empty(t, buf.String())
	})

	t.Run("slice of structs", func(t *testing.T) {
		t.Parallel()

		exporter := xcsv.NewCSVExporter()

		var buf bytes.Buffer

		now := time.Now()
		data := []TestPerson{
			{Name: "Alice", Age: 30, Active: true, JoinedAt: now},
			{Name: "Bob", Age: 25, Active: false, JoinedAt: now},
		}

		err := exporter.Export(&buf, data)
		require.NoError(t, err)

		result := buf.String()
		assert.Contains(t, result, "Name")
		assert.Contains(t, result, "Age")
		assert.Contains(t, result, "Active")
		assert.Contains(t, result, "JoinedAt")
		assert.Contains(t, result, "Alice")
		assert.Contains(t, result, "30")
		assert.Contains(t, result, "true")
		assert.Contains(t, result, "Bob")
		assert.Contains(t, result, "25")
		assert.Contains(t, result, "false")
	})

	t.Run("slice of maps", func(t *testing.T) {
		t.Parallel()

		exporter := xcsv.NewCSVExporter()

		var buf bytes.Buffer

		data := []map[string]any{
			{"name": "Alice", "age": 30},
			{"name": "Bob", "age": 25},
		}

		err := exporter.Export(&buf, data)
		require.NoError(t, err)

		result := buf.String()
		assert.Contains(t, result, "name")
		assert.Contains(t, result, "age")
		assert.Contains(t, result, "Alice")
		assert.Contains(t, result, "30")
		assert.Contains(t, result, "Bob")
		assert.Contains(t, result, "25")
	})

	t.Run("slice of primitive types", func(t *testing.T) {
		t.Parallel()

		exporter := xcsv.NewCSVExporter()

		var buf bytes.Buffer

		data := []int{1, 2, 3, 4, 5}

		err := exporter.Export(&buf, data)
		require.NoError(t, err)

		result := buf.String()
		assert.Equal(t, "1\n2\n3\n4\n5\n", result)
	})

	t.Run("custom delimiter", func(t *testing.T) {
		t.Parallel()

		exporter := xcsv.NewCSVExporter(xcsv.WithDelimiter(';'))

		var buf bytes.Buffer

		data := []map[string]any{
			{"name": "Alice", "age": 30},
			{"name": "Bob", "age": 25},
		}

		err := exporter.Export(&buf, data)
		require.NoError(t, err)

		result := buf.String()
		assert.Equal(t, "age;name\n30;Alice\n25;Bob\n", result)
	})

	t.Run("CRLF line endings", func(t *testing.T) {
		t.Parallel()

		exporter := xcsv.NewCSVExporter(xcsv.WithCRLF(true))

		var buf bytes.Buffer

		data := []string{"a", "b", "c"}

		err := exporter.Export(&buf, data)
		require.NoError(t, err)

		result := buf.String()
		assert.Contains(t, result, "\r\n")
	})

	t.Run("null values in maps", func(t *testing.T) {
		t.Parallel()

		exporter := xcsv.NewCSVExporter()

		var buf bytes.Buffer

		data := []map[string]any{
			{"name": "Alice", "age": nil},
			{"name": nil, "age": 25},
		}

		err := exporter.Export(&buf, data)
		require.NoError(t, err)

		result := buf.String()
		assert.Contains(t, result, "Alice")
		assert.Contains(t, result, "25")
		assert.Contains(t, result, "name")
		assert.Contains(t, result, "age")
	})

	t.Run("special characters", func(t *testing.T) {
		t.Parallel()

		exporter := xcsv.NewCSVExporter()

		var buf bytes.Buffer

		data := []map[string]any{
			{"field": "contains,comma"},
			{"field": "contains\"quote"},
			{"field": "contains\nnewline"},
		}

		err := exporter.Export(&buf, data)
		require.NoError(t, err)

		result := buf.String()
		assert.Contains(t, result, "\"contains,comma\"")
		assert.Contains(t, result, "\"contains\"\"quote\"")
		assert.Contains(t, result, "\"contains\nnewline\"")
	})

	t.Run("unicode characters", func(t *testing.T) {
		t.Parallel()

		exporter := xcsv.NewCSVExporter()

		var buf bytes.Buffer

		data := []string{"🌟", "世界", "über"}

		err := exporter.Export(&buf, data)
		require.NoError(t, err)

		result := buf.String()
		assert.Contains(t, result, "🌟")
		assert.Contains(t, result, "世界")
		assert.Contains(t, result, "über")
	})

	t.Run("private fields in struct", func(t *testing.T) {
		t.Parallel()

		exporter := xcsv.NewCSVExporter()

		var buf bytes.Buffer

		data := []PrivateFieldStruct{
			{Public: "visible", private: "hidden"},
		}

		err := exporter.Export(&buf, data)
		require.NoError(t, err)

		result := buf.String()
		assert.Contains(t, result, "Public")
		assert.NotContains(t, result, "private")
		assert.Contains(t, result, "visible")
		assert.NotContains(t, result, "hidden")
	})

	t.Run("missing map fields", func(t *testing.T) {
		t.Parallel()

		exporter := xcsv.NewCSVExporter()

		var buf bytes.Buffer

		data := []map[string]any{
			{"name": "Alice", "age": 30},
			{"name": "Bob"}, // missing age
		}

		err := exporter.Export(&buf, data)
		require.NoError(t, err)

		result := buf.String()
		assert.Contains(t, result, "Alice")
		assert.Contains(t, result, "30")
		assert.Contains(t, result, "Bob")
	})

	t.Run("empty map slice", func(t *testing.T) {
		t.Parallel()

		exporter := xcsv.NewCSVExporter()

		var buf bytes.Buffer

		data := []map[string]any{}

		err := exporter.Export(&buf, data)
		require.NoError(t, err)

		assert.Empty(t, buf.String())
	})
}
