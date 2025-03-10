// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package xjson_test

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/polyclient/polyclient/pkg/dataexchange/xjson"
	"github.com/stretchr/testify/assert"
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

func TestNewJsonExporter(t *testing.T) {
	t.Parallel()

	t.Run("default configuration", func(t *testing.T) {
		t.Parallel()

		exporter := xjson.NewJSONExporter()
		assert.NotNil(t, exporter)
	})

	t.Run("custom configuration", func(t *testing.T) {
		t.Parallel()

		exporter := xjson.NewJSONExporter(
			xjson.WithIndentString("\t"),
			xjson.WithEscapeHTML(false),
		)
		assert.NotNil(t, exporter)
	})
}

func TestExport(t *testing.T) {
	t.Parallel()

	t.Run("nil writer", func(t *testing.T) {
		t.Parallel()

		exporter := xjson.NewJSONExporter()
		err := exporter.Export(nil, []string{"test"})
		assert.Error(t, err)
	})

	t.Run("empty slice", func(t *testing.T) {
		t.Parallel()

		exporter := xjson.NewJSONExporter()

		var buf bytes.Buffer

		err := exporter.Export(&buf, []string{})
		assert.NoError(t, err)
		assert.Equal(t, "[]\n", buf.String())
	})

	t.Run("slice of structs", func(t *testing.T) {
		t.Parallel()

		exporter := xjson.NewJSONExporter()

		var buf bytes.Buffer

		now := time.Date(2025, 3, 14, 15, 9, 26, 0, time.UTC)
		data := []TestPerson{
			{Name: "Alice", Age: 30, Active: true, JoinedAt: now},
			{Name: "Bob", Age: 25, Active: false, JoinedAt: now},
		}

		err := exporter.Export(&buf, data)
		assert.NoError(t, err)

		var result []TestPerson
		err = json.Unmarshal(buf.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, data, result)
	})

	t.Run("slice of maps", func(t *testing.T) {
		t.Parallel()

		exporter := xjson.NewJSONExporter()

		var buf bytes.Buffer

		data := []map[string]any{
			{"name": "Alice", "age": float64(30)},
			{"name": "Bob", "age": float64(25)},
		}

		err := exporter.Export(&buf, data)
		assert.NoError(t, err)

		var result []map[string]any
		err = json.Unmarshal(buf.Bytes(), &result)
		assert.NoError(t, err)
		assert.Len(t, result, 2)

		foundAlice := false
		foundBob := false

		for _, m := range result {
			if m["name"] == "Alice" && m["age"] == float64(30) {
				foundAlice = true
			}

			if m["name"] == "Bob" && m["age"] == float64(25) {
				foundBob = true
			}
		}

		assert.True(t, foundAlice, "Alice's data not found")
		assert.True(t, foundBob, "Bob's data not found")
	})

	t.Run("null values in maps", func(t *testing.T) {
		t.Parallel()

		exporter := xjson.NewJSONExporter()

		var buf bytes.Buffer

		data := []map[string]any{
			{"name": "Alice", "age": nil},
			{"name": nil, "age": float64(25)},
		}

		err := exporter.Export(&buf, data)
		assert.NoError(t, err)

		var result []map[string]any
		err = json.Unmarshal(buf.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, data, result)
	})

	t.Run("html escaping enabled", func(t *testing.T) {
		t.Parallel()

		exporter := xjson.NewJSONExporter(xjson.WithEscapeHTML(true))

		var buf bytes.Buffer

		data := map[string]string{
			"html": "<script>alert('xss')</script>",
			"text": "a & b < c > d",
		}

		err := exporter.Export(&buf, data)
		assert.NoError(t, err)

		result := buf.String()
		assert.Contains(t, result, `\u003cscript\u003e`)
		assert.Contains(t, result, `\u003c`)
		assert.Contains(t, result, `\u003e`)
		assert.Contains(t, result, `\u0026`)
	})

	t.Run("html escaping disabled", func(t *testing.T) {
		t.Parallel()

		exporter := xjson.NewJSONExporter(xjson.WithEscapeHTML(false))

		var buf bytes.Buffer

		data := map[string]string{
			"html": "<script>alert('xss')</script>",
			"text": "a & b < c > d",
		}

		err := exporter.Export(&buf, data)
		assert.NoError(t, err)

		result := buf.String()
		assert.Contains(t, result, "<script>")
		assert.Contains(t, result, "&")
		assert.Contains(t, result, "<")
		assert.Contains(t, result, ">")
	})

	t.Run("custom indentation", func(t *testing.T) {
		t.Parallel()

		exporter := xjson.NewJSONExporter(xjson.WithIndentString("\t"))

		var buf bytes.Buffer

		data := map[string]string{"key": "value"}

		err := exporter.Export(&buf, data)
		assert.NoError(t, err)

		result := buf.String()
		assert.Contains(t, result, "\t\"key\":")
	})

	t.Run("unicode characters", func(t *testing.T) {
		t.Parallel()

		exporter := xjson.NewJSONExporter()

		var buf bytes.Buffer

		data := []string{"🌟", "世界", "über"}

		err := exporter.Export(&buf, data)
		assert.NoError(t, err)

		var result []string
		err = json.Unmarshal(buf.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, data, result)
	})

	t.Run("private fields in struct", func(t *testing.T) {
		t.Parallel()

		exporter := xjson.NewJSONExporter()

		var buf bytes.Buffer

		data := []PrivateFieldStruct{
			{Public: "visible", private: "hidden"},
		}

		err := exporter.Export(&buf, data)
		assert.NoError(t, err)

		var result []map[string]any
		err = json.Unmarshal(buf.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, "visible", result[0]["Public"])
		assert.NotContains(t, result[0], "private")
	})
}
