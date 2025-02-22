// Copyright (C) 2025 Juan Mesa and contributors
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3
// as published by the Free Software Foundation, with the Runtime
// Library Exception. See the COPYING.RUNTIME file for details.

package exporter

import (
	"encoding/csv"
	"io"

	"github.com/jszwec/csvutil"
)

// FormatTSV formats and writes the provided data to the writer in TSV format.
func FormatTSV(w io.Writer, data any) error {
	writer := csv.NewWriter(w)
	writer.Comma = '\t'
	defer writer.Flush()

	return csvutil.NewEncoder(writer).Encode(data)
}
