// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package xhtml

import (
	"errors"
	"fmt"
	"html"
	"html/template"
	"io"
	"reflect"
	"time"

	"github.com/polyclient/polyclient/pkg/stringify"
)

// HtmlExporter is a data exporter for HTML format.
type HtmlExporter struct {
	// DateFormat is the format for date fields (default time.RFC3339).
	DateFormat string
	// UseCss is whether to use default styles in the HTML output (default true).
	UseCss bool
	// template is the template for the HTML output.
	template *template.Template
}

// HtmlExporterOption defines a functional option for configuring HTML output.
type HtmlExporterOption func(*HtmlExporter)

// WithDateFormat sets a custom date format for date fields.
func WithDateFormat(format string) HtmlExporterOption {
	return func(opts *HtmlExporter) {
		opts.DateFormat = format
	}
}

// WithUseCss sets whether to use default styles in the HTML output.
func WithUseCss(useCss bool) HtmlExporterOption {
	return func(opts *HtmlExporter) {
		opts.UseCss = useCss
	}
}

// NewHtmlExporter creates a new HtmlExporter with the specified options.
func NewHtmlExporter(opts ...HtmlExporterOption) *HtmlExporter {
	ex := &HtmlExporter{
		DateFormat: time.RFC3339,
		UseCss:     true,
		template:   GetTemplate(),
	}

	for _, opt := range opts {
		opt(ex)
	}

	return ex
}

// Export writes the provided slice data to the given writer in HTML format.
// The data must be a slice of one of:
// - structs (writes headers from exported fields)
// - maps[string]any (writes headers from first row's keys)
// Returns an error if data is not a slice or if writing fails.
func (ex *HtmlExporter) Export(w io.Writer, data any) error {
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

	converted := make([]any, v.Len())
	for i := range v.Len() {
		converted[i] = v.Index(i).Interface()
	}

	parsedData, err := ex.formatSlice(converted)
	if err != nil {
		return err
	}

	parsedData.UseCss = ex.UseCss

	return ex.template.Execute(w, parsedData)
}

// formatSlice formats the data for the HTML template.
func (ex *HtmlExporter) formatSlice(data []any) (*HtmlTemplateData, error) {
	switch first := data[0].(type) {
	case map[string]any:
		return ex.formatMapSlice(data)
	default:
		if reflect.TypeOf(first).Kind() == reflect.Struct {
			return ex.formatStructSlice(data)
		}

		return ex.formatSingleColumnSlice(data)
	}
}

// formatMapSlice formats a slice of maps for HTML output.
func (ex *HtmlExporter) formatMapSlice(data []any) (*HtmlTemplateData, error) {
	parsed := &HtmlTemplateData{
		Headers: make([]string, 0),
		Rows:    make([][]string, 0, len(data)),
	}

	first, ok := data[0].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("first element is not a map: %T", data[0])
	}

	for header := range first {
		parsed.Headers = append(parsed.Headers, header)
	}

	for _, item := range data {
		record, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("item is not a map: %T", item)
		}

		row := make([]string, len(parsed.Headers))
		for i, header := range parsed.Headers {
			row[i] = stringify.Stringify(record[header],
				stringify.WithDateFormat(ex.DateFormat),
				stringify.WithCustomFormatter(sanitizeHtml),
			)
		}

		parsed.Rows = append(parsed.Rows, row)
	}

	return parsed, nil
}

// formatStructSlice formats a slice of structs for HTML output.
func (ex *HtmlExporter) formatStructSlice(data []any) (*HtmlTemplateData, error) {
	parsed := &HtmlTemplateData{
		Headers: make([]string, 0),
		Rows:    make([][]string, 0, len(data)),
	}

	first := reflect.TypeOf(data[0])
	for i := 0; i < first.NumField(); i++ {
		field := first.Field(i)
		if field.PkgPath == "" {
			parsed.Headers = append(parsed.Headers, field.Name)
		}
	}

	for _, item := range data {
		value := reflect.ValueOf(item)
		row := make([]string, len(parsed.Headers))

		for i, header := range parsed.Headers {
			field := value.FieldByName(header)
			if field.IsValid() && field.CanInterface() {
				row[i] = stringify.Stringify(field.Interface(),
					stringify.WithDateFormat(ex.DateFormat),
					stringify.WithCustomFormatter(sanitizeHtml),
				)
			} else {
				row[i] = ""
			}
		}

		parsed.Rows = append(parsed.Rows, row)
	}

	return parsed, nil
}

// formatSingleColumnSlice writes `[]any` as a single-column HTML.
func (ex *HtmlExporter) formatSingleColumnSlice(data []any) (*HtmlTemplateData, error) {
	parsed := &HtmlTemplateData{
		Headers: []string{"Value"},
		Rows:    make([][]string, 0, len(data)),
	}

	for _, item := range data {
		parsed.Rows = append(parsed.Rows, []string{
			stringify.Stringify(item,
				stringify.WithDateFormat(ex.DateFormat),
				stringify.WithCustomFormatter(sanitizeHtml),
			),
		})
	}

	return parsed, nil
}

// sanitizeHtml escapes HTML characters in the provided string.
func sanitizeHtml(v string) string {
	return html.EscapeString(v)
}
