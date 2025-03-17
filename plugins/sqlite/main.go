// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/polyclient/polyclient/bindings/polyclient/sql/query"
	"github.com/polyclient/polyclient/plugins/sqlite/internal"
)

func init() {
	query.Exports.Execute = internal.NewExecute
	// query.Exports.ExecuteBatch = NewExecuteBatch
	// query.Exports.ExecutePrepared = NewExecutePrepared
}

func main() {}
