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
	"sort"
	"time"

	"github.com/polyclient/polyclient/pkg/stringify"
)

// HTMLExporter is a data exporter for HTML format.
type HTMLExporter struct {
	// DateFormat is the format for date fields (default time.RFC3339).
	DateFormat string
	// UseCSS is whether to use default CSS styles in the HTML output (default true).
	UseCSS bool
	// template is the template for the HTML output.
	template *template.Template
}

// HTMLExporterOption defines a functional option for configuring HTML output.
type HTMLExporterOption func(*HTMLExporter)

// WithDateFormat sets a custom date format for date fields.
func WithDateFormat(format string) HTMLExporterOption {
	return func(opts *HTMLExporter) {
		opts.DateFormat = format
	}
}

// WithUseCSS sets whether to use default styles in the HTML output.
func WithUseCSS(useCSS bool) HTMLExporterOption {
	return func(opts *HTMLExporter) {
		opts.UseCSS = useCSS
	}
}

// NewHTMLExporter creates a new HtmlExporter with the specified options.
func NewHTMLExporter(opts ...HTMLExporterOption) *HTMLExporter {
	ex := &HTMLExporter{
		DateFormat: time.RFC3339,
		UseCSS:     true,
		template:   GetTemplate(),
	}

	for _, opt := range opts {
		opt(ex)
	}

	return ex
}

// Export writes a slice to HTML, supporting primitive types, structs, and maps.
func (ex *HTMLExporter) Export(w io.Writer, data any) error {
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

	parsedData, err := ex.formatSlice(v)
	if err != nil {
		return err
	}

	parsedData.UseCSS = ex.UseCSS

	return ex.template.Execute(w, parsedData)
}

// formatSlice formats the data for the HTML template.
func (ex *HTMLExporter) formatSlice(v reflect.Value) (*HTMLTemplateData, error) {
	first := v.Index(0).Interface()

	switch first.(type) {
	case map[string]any:
		return ex.formatMapSlice(v)
	default:
		if v.Index(0).Kind() == reflect.Struct {
			return ex.formatStructSlice(v)
		}

		return ex.formatSingleColumnSlice(v)
	}
}

// formatMapSlice formats a slice of maps for HTML output.
func (ex *HTMLExporter) formatMapSlice(v reflect.Value) (*HTMLTemplateData, error) {
	parsed := &HTMLTemplateData{
		Headers: []string{},
		Rows:    make([][]string, 0, v.Len()),
	}

	first := v.Index(0).Interface().(map[string]any)
	for header := range first {
		parsed.Headers = append(parsed.Headers, header)
	}

	sort.Strings(parsed.Headers) // Ensure consistent column ordering

	for i := range v.Len() {
		record := v.Index(i).Interface().(map[string]any)
		row := make([]string, len(parsed.Headers))

		for j, header := range parsed.Headers {
			row[j] = stringify.Stringify(record[header],
				stringify.WithDateFormat(ex.DateFormat),
				stringify.WithCustomFormatter(sanitizeHTML),
			)
		}

		parsed.Rows = append(parsed.Rows, row)
	}

	return parsed, nil
}

// formatStructSlice formats a slice of structs for HTML output.
func (ex *HTMLExporter) formatStructSlice(v reflect.Value) (*HTMLTemplateData, error) {
	parsed := &HTMLTemplateData{
		Headers: []string{},
		Rows:    make([][]string, 0, v.Len()),
	}

	typeOfStruct := v.Index(0).Type()
	for i := range typeOfStruct.NumField() {
		field := typeOfStruct.Field(i)
		if field.PkgPath == "" {
			parsed.Headers = append(parsed.Headers, field.Name)
		}
	}

	for i := range v.Len() {
		value := v.Index(i)
		row := make([]string, len(parsed.Headers))

		for j, header := range parsed.Headers {
			field := value.FieldByName(header)
			if field.IsValid() && field.CanInterface() {
				row[j] = stringify.Stringify(field.Interface(),
					stringify.WithDateFormat(ex.DateFormat),
					stringify.WithCustomFormatter(sanitizeHTML),
				)
			} else {
				row[j] = ""
			}
		}

		parsed.Rows = append(parsed.Rows, row)
	}

	return parsed, nil
}

// formatSingleColumnSlice writes `[]any` as a single-column HTML.
func (ex *HTMLExporter) formatSingleColumnSlice(v reflect.Value) (*HTMLTemplateData, error) {
	parsed := &HTMLTemplateData{
		Headers: []string{"Value"},
		Rows:    make([][]string, 0, v.Len()),
	}

	for i := range v.Len() {
		parsed.Rows = append(parsed.Rows, []string{
			stringify.Stringify(v.Index(i).Interface(),
				stringify.WithDateFormat(ex.DateFormat),
				stringify.WithCustomFormatter(sanitizeHTML),
			),
		})
	}

	return parsed, nil
}

// sanitizeHTML escapes HTML characters in the provided string.
func sanitizeHTML(v string) string {
	return html.EscapeString(v)
}
