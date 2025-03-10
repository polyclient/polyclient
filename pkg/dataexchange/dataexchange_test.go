// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package dataexchange_test

import (
	"encoding/json"
	"testing"

	"github.com/polyclient/polyclient/pkg/dataexchange"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSupportedExportFormats(t *testing.T) {
	t.Parallel()

	expectedFormats := []dataexchange.Format{
		dataexchange.FormatCsv, dataexchange.FormatTsv, dataexchange.FormatJson, dataexchange.FormatHtml,
	}

	supportedFormats := dataexchange.GetSupportedExportFormats()
	assert.ElementsMatch(t, expectedFormats, supportedFormats)
}

func TestGetExporterRegistryEntry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		format   dataexchange.Format
		exists   bool
		mimeType string
		fileExt  string
	}{
		{"csv", dataexchange.FormatCsv, true, "text/csv", "csv"},
		{"tsv", dataexchange.FormatTsv, true, "text/tab-separated-values", "tsv"},
		{"json", dataexchange.FormatJson, true, "application/json", "json"},
		{"html", dataexchange.FormatHtml, true, "text/html", "html"},
		{"Invalid", "invalid_format", false, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			entry, exists := dataexchange.GetRegistryEntry(tt.format)
			assert.Equal(t, tt.exists, exists)

			if exists {
				assert.Equal(t, tt.mimeType, entry.MIMEType)
				assert.Equal(t, tt.fileExt, entry.FileExt)
			}
		})
	}
}

func TestParseDataFromBytes_ValidJSON(t *testing.T) {
	t.Parallel()

	data := []byte(`[{"key": "value"}, {"key": "value2"}]`)

	var expected []map[string]any

	require.NoError(t, json.Unmarshal(data, &expected))

	parsed, err := dataexchange.ParseDataFromBytes[[]map[string]any](data, dataexchange.FormatJson)
	require.NoError(t, err)
	assert.Equal(t, expected, parsed)
}

func TestParseDataFromBytes_InvalidJSON(t *testing.T) {
	t.Parallel()

	data := []byte(`invalid json`)

	var zeroValue []map[string]any

	parsed, err := dataexchange.ParseDataFromBytes[[]map[string]any](data, dataexchange.FormatJson)
	assert.Error(t, err)
	assert.Equal(t, zeroValue, parsed)
}

func TestParseDataFromBytes_UnsupportedFormat(t *testing.T) {
	t.Parallel()

	data := []byte(`[{"key": "value"}]`)

	var zeroValue []map[string]any

	parsed, err := dataexchange.ParseDataFromBytes[[]map[string]any](data, "unsupported")
	assert.Error(t, err)
	assert.Equal(t, zeroValue, parsed)
}
