// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

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
