// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package exporter

import (
	"fmt"
	"io"
)

// Exporter represents the data format exporter.
type DataExporter struct {
	format Format
	writer io.Writer
}

// DataExporterOptions represents the options for creating a DataExporter.
type DataExporterOptions struct {
	Format string
	Writer io.Writer
}

// NewDataExporter creates a new DataExporter instance with the given options.
func NewDataExporter(opts DataExporterOptions) *DataExporter {
	return &DataExporter{
		format: Format(opts.Format),
		writer: opts.Writer,
	}
}

// Exports the given data to the specified format using the configured writer.
func (de *DataExporter) Export(data any) error {
	formatter, ok := formatRegistry[de.format]
	if !ok {
		return fmt.Errorf("format %s is not supported", de.format)
	}

	return formatter.Formatter(de.writer, data)
}
