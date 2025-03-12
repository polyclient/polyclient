// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package datapipe_test

import (
	"encoding/json"
	"testing"

	"github.com/polyclient/polyclient/internal/datapipe"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSupportedExportFormats(t *testing.T) {
	t.Parallel()

	expectedFormats := []datapipe.Format{
		datapipe.FormatCSV, datapipe.FormatTSV, datapipe.FormatJSON, datapipe.FormatHTML,
	}

	supportedFormats := datapipe.GetAvailableExportFormats()
	assert.ElementsMatch(t, expectedFormats, supportedFormats)
}

func TestGetExporterRegistryEntry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		format   datapipe.Format
		exists   bool
		mimeType string
		fileExt  string
	}{
		{"csv", datapipe.FormatCSV, true, "text/csv", "csv"},
		{"tsv", datapipe.FormatTSV, true, "text/tab-separated-values", "tsv"},
		{"json", datapipe.FormatJSON, true, "application/json", "json"},
		{"html", datapipe.FormatHTML, true, "text/html", "html"},
		{"Invalid", "invalid_format", false, "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			entry, exists := datapipe.GetRegistryEntry(tt.format)
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

	parsed, err := datapipe.ParseDataFromBytes[[]map[string]any](data, datapipe.FormatJSON)
	require.NoError(t, err)

	assert.Equal(t, expected, parsed)
}

func TestParseDataFromBytes_InvalidJSON(t *testing.T) {
	t.Parallel()

	data := []byte(`invalid json`)

	var zeroValue []map[string]any

	parsed, err := datapipe.ParseDataFromBytes[[]map[string]any](data, datapipe.FormatJSON)
	require.Error(t, err)

	assert.Equal(t, zeroValue, parsed)
}

func TestParseDataFromBytes_UnsupportedFormat(t *testing.T) {
	t.Parallel()

	data := []byte(`[{"key": "value"}]`)

	var zeroValue []map[string]any

	parsed, err := datapipe.ParseDataFromBytes[[]map[string]any](data, "unsupported")
	require.Error(t, err)

	assert.Equal(t, zeroValue, parsed)
}
