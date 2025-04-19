// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package db

import "context"

// QueryExecutor defines the capability to execute queries.
type QueryExecutor interface {
	// Execute runs a query against the connection.
	// The 'query' parameter is driver-specific (e.g., SQL string, BSON document, Cypher).
	// Use driver-specific query builder objects or simple types like string/map.
	Execute(ctx context.Context, query string) (Result, error)
}
