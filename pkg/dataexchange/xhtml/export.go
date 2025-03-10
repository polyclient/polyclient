// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package xhtml

import (
	"fmt"
	"html"
	"html/template"
	"io"
	"reflect"
	"time"

	"github.com/polyclient/polyclient/pkg/stringify"
	"github.com/samber/lo"
)

// HtmlExporter is a data exporter for HTML format.
type HtmlExporter struct {
	// DateFormat is the format for date fields (default time.RFC3339).
	DateFormat string
	// UseCss is whether to use default styles in the HTML output (default true).
	UseCss bool

	// Template is the template for the HTML output.
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

// NewHtmlExporter creates a new HtmlExporter.
func NewHtmlExporter(opts ...HtmlExporterOption) *HtmlExporter {
	ex := &HtmlExporter{
		DateFormat: time.RFC3339,
		template:   GetTemplate(),
		UseCss:     true,
	}

	for _, opt := range opts {
		opt(ex)
	}

	return ex
}

// Export formats and writes the provided data to the writer in HTML format.
func (ex *HtmlExporter) Export(w io.Writer, data any) error {
	v := reflect.ValueOf(data)

	if v.Kind() != reflect.Slice {
		return fmt.Errorf("expected a slice, got %T", data)
	}

	if v.Len() == 0 {
		return nil
	}

	var converted []any

	elemType := v.Type().Elem()
	switch elemType.Kind() {
	case reflect.Struct:
		converted = make([]any, v.Len())
		for i := range v.Len() {
			converted[i] = v.Index(i).Interface()
		}
	case reflect.Interface:
		converted = data.([]any)
	default:
		return fmt.Errorf("unsupported data type: %T", data)
	}

	parsedData, err := ex.formatSlice(converted)
	if err != nil {
		return err
	}

	return ex.template.Execute(w, parsedData)
}

// formatDataForTemplate formats the data for the HTML template.
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

// formatMapSlice writes `[]map[string]any` as a multi-column HTML table.
func (ex *HtmlExporter) formatMapSlice(data []any) (*HtmlTemplateData, error) {
	parsed := &HtmlTemplateData{
		Headers: []string{},
		Rows:    [][]string{},
	}

	first, ok := data[0].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("first element is not a map: %T", data[0])
	}

	parsed.Headers = lo.Uniq(lo.Keys(first))

	for _, item := range data {
		record, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("item is not a map: %T", item)
		}

		row := lo.Map(parsed.Headers, func(header string, _ int) string {
			return stringify.Stringify(record[header], stringify.WithCustomFormatter(sanitizeHtml))
		})

		parsed.Rows = append(parsed.Rows, row)
	}

	return parsed, nil
}

// formatStructSlice writes `[]struct` as a multi-column HTML table.
func (ex *HtmlExporter) formatStructSlice(data []any) (*HtmlTemplateData, error) {
	parsed := &HtmlTemplateData{
		Headers: []string{},
		Rows:    [][]string{},
	}

	first := reflect.TypeOf(data[0])
	for i := range first.NumField() {
		field := first.Field(i)

		if field.PkgPath == "" {
			parsed.Headers = append(parsed.Headers, field.Name)
		}
	}

	parsed.Rows = lo.Map(data, func(item any, _ int) []string {
		value := reflect.ValueOf(item)

		return lo.Map(parsed.Headers, func(header string, _ int) string {
			field := value.FieldByName(header)

			if field.IsValid() && field.CanInterface() {
				return stringify.Stringify(field.Interface(), stringify.WithCustomFormatter(sanitizeHtml))
			}

			return ""
		})
	})

	return parsed, nil
}

// formatSingleColumnSlice writes `[]any` as a single-column HTML table.
func (ex *HtmlExporter) formatSingleColumnSlice(data []any) (*HtmlTemplateData, error) {
	parsed := &HtmlTemplateData{
		Headers: []string{},
		Rows:    [][]string{},
	}

	parsed.Rows = lo.Map(data, func(item any, _ int) []string {
		return []string{stringify.Stringify(item, stringify.WithCustomFormatter(sanitizeHtml))}
	})

	return parsed, nil
}

// sanitizeHtml escapes HTML characters in the provided string.
func sanitizeHtml(v string) string {
	return html.EscapeString(v)
}
