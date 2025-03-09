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
	"github.com/stretchr/testify/require"
)

func TestJsonExporter_Defaults(t *testing.T) {
	t.Parallel()

	exporter := xjson.NewJsonExporter()

	var buf bytes.Buffer

	data := map[string]any{"key": "value"}
	err := exporter.Export(&buf, data)
	require.NoError(t, err)

	var result map[string]any

	require.NoError(t, json.Unmarshal(buf.Bytes(), &result))
	assert.Equal(t, data, result)
}

func TestJsonExporter_Indentation(t *testing.T) {
	t.Parallel()

	exporter := xjson.NewJsonExporter(xjson.WithIndentString("\t"))

	var buf bytes.Buffer

	data := map[string]any{"key": "value"}
	require.NoError(t, exporter.Export(&buf, data))

	assert.Contains(t, buf.String(), "\n\t") // Check if tab indentation is applied
}

func TestJsonExporter_EscapeHTML(t *testing.T) {
	t.Parallel()

	data := map[string]any{"html": "<script>alert('xss')</script>"}

	exporter := xjson.NewJsonExporter()

	var buf bytes.Buffer

	require.NoError(t, exporter.Export(&buf, data))
	assert.Contains(t, buf.String(), "\\u003cscript\\u003e") // Default should escape HTML

	exporter = xjson.NewJsonExporter(xjson.WithEscapeHTML(false))

	buf.Reset()
	require.NoError(t, exporter.Export(&buf, data))
	assert.Contains(t, buf.String(), "<script>") // HTML should not be escaped
}

func TestJsonExporter_VariousDataTypes(t *testing.T) {
	t.Parallel()

	exporter := xjson.NewJsonExporter()

	var buf bytes.Buffer

	timeVal := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	data := map[string]any{
		"string": "hello",
		"int":    42,
		"float":  3.14,
		"bool":   true,
		"nil":    nil,
		"time":   timeVal,
	}
	require.NoError(t, exporter.Export(&buf, data))

	var result map[string]any

	require.NoError(t, json.Unmarshal(buf.Bytes(), &result))
	assert.Equal(t, "hello", result["string"])
	assert.Equal(t, float64(42), result["int"]) // JSON encodes numbers as float64
	assert.Equal(t, 3.14, result["float"])
	assert.Equal(t, true, result["bool"])
	assert.Nil(t, result["nil"])
	assert.Equal(t, timeVal.Format(time.RFC3339), result["time"])
}
