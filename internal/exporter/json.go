// Copyright (C) 2025 Juan Mesa and contributors
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3
// as published by the Free Software Foundation, with the Runtime
// Library Exception. See the COPYING.RUNTIME file for details.

package exporter

import (
	"encoding/json"
	"io"
)

// formatJSON formats and writes the provided data to the writer in JSON format.
// It automatically adds indentation of 2 spaces to make the output more readable.
// The data parameter can be of any type that is JSON-encodable.
// Returns an error if JSON encoding fails.
func FormatJSON(w io.Writer, data any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	return enc.Encode(data)
}
