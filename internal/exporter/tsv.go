// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

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
