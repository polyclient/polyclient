// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package db

// Config represents the configuration for opening a database connection.
// Using a map allows flexibility, but drivers should document expected keys.
type Config map[string]any
