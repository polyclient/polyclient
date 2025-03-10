// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package xhtml_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/polyclient/polyclient/pkg/dataexchange/xhtml"
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

func TestNewHtmlExporter(t *testing.T) {
	t.Parallel()

	t.Run("default configuration", func(t *testing.T) {
		t.Parallel()

		exporter := xhtml.NewHtmlExporter()
		assert.NotNil(t, exporter)
	})

	t.Run("custom configuration", func(t *testing.T) {
		t.Parallel()

		exporter := xhtml.NewHtmlExporter(
			xhtml.WithDateFormat("2006-01-02"),
			xhtml.WithUseCss(false),
		)
		assert.NotNil(t, exporter)
	})
}

func TestExport(t *testing.T) {
	t.Parallel()

	t.Run("invalid input - non-slice", func(t *testing.T) {
		t.Parallel()

		exporter := xhtml.NewHtmlExporter()

		var buf bytes.Buffer
		err := exporter.Export(&buf, "not a slice")
		assert.Error(t, err)
	})

	t.Run("empty slice", func(t *testing.T) {
		t.Parallel()

		exporter := xhtml.NewHtmlExporter()

		var buf bytes.Buffer
		err := exporter.Export(&buf, []string{})
		assert.NoError(t, err)
		assert.Empty(t, buf.String())
	})

	t.Run("slice of structs", func(t *testing.T) {
		t.Parallel()

		exporter := xhtml.NewHtmlExporter()

		var buf bytes.Buffer

		now := time.Now()
		data := []TestPerson{
			{Name: "Alice", Age: 30, Active: true, JoinedAt: now},
			{Name: "Bob", Age: 25, Active: false, JoinedAt: now},
		}

		err := exporter.Export(&buf, data)
		assert.NoError(t, err)

		result := buf.String()
		assert.Contains(t, result, "<th>Name</th>")
		assert.Contains(t, result, "<th>Age</th>")
		assert.Contains(t, result, "<th>Active</th>")
		assert.Contains(t, result, "<th>JoinedAt</th>")
		assert.Contains(t, result, "<td>Alice</td>")
		assert.Contains(t, result, "<td>30</td>")
		assert.Contains(t, result, "<td>true</td>")
	})

	t.Run("slice of maps", func(t *testing.T) {
		t.Parallel()

		exporter := xhtml.NewHtmlExporter()

		var buf bytes.Buffer

		data := []map[string]any{
			{"name": "Alice", "age": 30},
			{"name": "Bob", "age": 25},
		}

		err := exporter.Export(&buf, data)
		assert.NoError(t, err)

		result := buf.String()
		assert.Contains(t, result, "<th>name</th>")
		assert.Contains(t, result, "<th>age</th>")
		assert.Contains(t, result, "<td>Alice</td>")
		assert.Contains(t, result, "<td>30</td>")
		assert.Contains(t, result, "<td>Bob</td>")
		assert.Contains(t, result, "<td>25</td>")
	})

	t.Run("null values in maps", func(t *testing.T) {
		t.Parallel()

		exporter := xhtml.NewHtmlExporter()

		var buf bytes.Buffer

		data := []map[string]any{
			{"name": "Alice", "age": nil},
			{"name": nil, "age": 25},
		}

		err := exporter.Export(&buf, data)
		assert.NoError(t, err)

		result := buf.String()
		assert.Contains(t, result, "<td>Alice</td>")
		assert.Contains(t, result, "<td>25</td>")
		assert.Contains(t, result, "<td></td>") // Empty cells for nil values
	})

	t.Run("html escaping enabled", func(t *testing.T) {
		t.Parallel()

		exporter := xhtml.NewHtmlExporter()

		var buf bytes.Buffer

		data := []map[string]any{
			{"field": "<script>alert('xss')</script>"},
			{"field": "Contains & and >"},
		}

		err := exporter.Export(&buf, data)
		assert.NoError(t, err)

		result := buf.String()
		assert.Contains(t, result, "&amp;lt;script&amp;gt;")
		assert.Contains(t, result, "Contains &amp;amp; and &amp;gt;")
		assert.NotContains(t, result, "<script>")
	})

	t.Run("unicode characters", func(t *testing.T) {
		t.Parallel()

		exporter := xhtml.NewHtmlExporter()

		var buf bytes.Buffer

		data := []map[string]any{
			{"field": "🌟"},
			{"field": "世界"},
			{"field": "über"},
		}

		err := exporter.Export(&buf, data)
		assert.NoError(t, err)

		result := buf.String()
		assert.Contains(t, result, "🌟")
		assert.Contains(t, result, "世界")
		assert.Contains(t, result, "über")
	})

	t.Run("private fields in struct", func(t *testing.T) {
		t.Parallel()

		exporter := xhtml.NewHtmlExporter()

		var buf bytes.Buffer

		data := []PrivateFieldStruct{
			{Public: "visible", private: "hidden"},
		}

		err := exporter.Export(&buf, data)
		assert.NoError(t, err)

		result := buf.String()
		assert.Contains(t, result, "<th>Public</th>")
		assert.NotContains(t, result, "private")
		assert.Contains(t, result, "<td>visible</td>")
		assert.NotContains(t, result, "hidden")
	})

	t.Run("css disabled", func(t *testing.T) {
		t.Parallel()

		exporter := xhtml.NewHtmlExporter(xhtml.WithUseCss(false))

		var buf bytes.Buffer

		data := []map[string]any{
			{"name": "Alice"},
		}

		err := exporter.Export(&buf, data)
		assert.NoError(t, err)

		result := buf.String()
		assert.NotContains(t, result, "<style>")
		assert.NotContains(t, result, "class=")
	})

	t.Run("custom date format", func(t *testing.T) {
		t.Parallel()

		exporter := xhtml.NewHtmlExporter(xhtml.WithDateFormat("2006-01-02"))

		var buf bytes.Buffer

		now := time.Date(2025, 3, 14, 15, 9, 26, 0, time.UTC)
		data := []TestPerson{
			{Name: "Alice", JoinedAt: now},
		}

		err := exporter.Export(&buf, data)
		assert.NoError(t, err)

		result := buf.String()
		assert.Contains(t, result, "2025-03-14")
		assert.NotContains(t, result, "15:09:26")
	})
}
