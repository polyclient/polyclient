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
	"github.com/stretchr/testify/require"
)

func TestNewHtmlExporter_Defaults(t *testing.T) {
	t.Parallel()

	ex := xhtml.NewHtmlExporter()
	assert.True(t, ex.UseCss)
	assert.Equal(t, time.RFC3339, ex.DateFormat)
}

func TestNewHtmlExporter_WithOptions(t *testing.T) {
	t.Parallel()

	ex := xhtml.NewHtmlExporter(
		xhtml.WithUseCss(false),
		xhtml.WithDateFormat("2006-01-02"),
	)

	assert.False(t, ex.UseCss)
	assert.Equal(t, "2006-01-02", ex.DateFormat)
}

func TestHtmlExporter_Export_SingleColumn(t *testing.T) {
	t.Parallel()

	data := []any{"foo", "bar", "baz"}

	var buf bytes.Buffer

	ex := xhtml.NewHtmlExporter()
	err := ex.Export(&buf, data)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "<td>foo</td>")
	assert.Contains(t, output, "<td>bar</td>")
	assert.Contains(t, output, "<td>baz</td>")
}

func TestHtmlExporter_Export_MapSlice(t *testing.T) {
	t.Parallel()

	data := []any{
		map[string]any{"name": "John", "age": 30},
		map[string]any{"name": "Jane", "age": 25},
	}

	var buf bytes.Buffer

	ex := xhtml.NewHtmlExporter()
	err := ex.Export(&buf, data)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "<th>name</th>")
	assert.Contains(t, output, "<th>age</th>")
	assert.Contains(t, output, "<td>John</td>")
	assert.Contains(t, output, "<td>30</td>")
	assert.Contains(t, output, "<td>Jane</td>")
	assert.Contains(t, output, "<td>25</td>")
}

func TestHtmlExporter_Export_StructSlice(t *testing.T) {
	t.Parallel()

	type Person struct {
		Name string
		Age  int
	}

	data := []Person{
		{Name: "Alice", Age: 28},
		{Name: "Bob", Age: 35},
	}

	var buf bytes.Buffer

	ex := xhtml.NewHtmlExporter()
	err := ex.Export(&buf, data)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "<th>Name</th>")
	assert.Contains(t, output, "<th>Age</th>")
	assert.Contains(t, output, "<td>Alice</td>")
	assert.Contains(t, output, "<td>28</td>")
	assert.Contains(t, output, "<td>Bob</td>")
	assert.Contains(t, output, "<td>35</td>")
}

func TestHtmlExporter_Export_EmptyInput(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	ex := xhtml.NewHtmlExporter()
	err := ex.Export(&buf, []any{})
	assert.NoError(t, err)
	assert.Empty(t, buf.String())
}

func TestHtmlExporter_Export_UnsupportedType(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	ex := xhtml.NewHtmlExporter()
	err := ex.Export(&buf, 123)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected a slice, got int")
}

func TestHtmlExporter_Export_InvalidMapSlice(t *testing.T) {
	t.Parallel()

	data := []any{
		map[string]any{"name": "Alice", "age": 30},
		"invalid", // Should cause an error
	}

	var buf bytes.Buffer

	ex := xhtml.NewHtmlExporter()
	err := ex.Export(&buf, data)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "item is not a map")
}
