// Copyright (C) 2025 Juan Mesa and contributors
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3
// as published by the Free Software Foundation, with the Runtime
// Library Exception. See the COPYING.RUNTIME file for details.

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
