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

// Exporter represents the data format exporter
type DataExporter struct {
	format Format
	output io.Writer
}

type DataExporterOptions struct {
	Format string
	Output io.Writer
}

func NewDataExporter(opts DataExporterOptions) *DataExporter {
	return &DataExporter{
		format: Format(opts.Format),
		output: opts.Output,
	}
}

func (de *DataExporter) Export(data any) error {
	formatter, ok := formatRegistry[de.format]
	if !ok {
		return fmt.Errorf("format %s is not supported", de.format)
	}

	return formatter.Formatter(de.output, data)
}
