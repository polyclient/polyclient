// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package database

import (
	"context"
	"time"
)

// Connection represents a database connection.
type Connection interface {
	// Query runs a query that returns rows.
	Query(ctx context.Context, query string, args ...any) (Rows, error)

	// Execute runs a query that doesn't return rows.
	Execute(ctx context.Context, query string, args ...any) (Result, error)

	// BeginTransaction starts a transaction.
	BeginTransaction(ctx context.Context, opts *TransactionOptions) (Transaction, error)

	// Ping checks if the connection is alive.
	Ping(ctx context.Context) error

	// Close closes the connection.
	Close()

	// ConnectionInfo returns information about the connection.
	ConnectionInfo() ConnectionInfo

	// Raw returns the underlying raw connection object.
	Raw() any
}

// ConnectionInfo represents information about a database connection.
type ConnectionInfo struct {
	DatabaseName  string
	ServerVersion string
	ServerType    string
	Host          string
	Port          int
	Username      string
	IsEncrypted   bool
	IsSSHEnabled  bool
	IsConnected   bool
	ConnectedAt   time.Time
	Stats         ConnectionStats
}

// ConnectionStats provides runtime statistics for a database connection.
type ConnectionStats struct {
	OpenTime         time.Duration
	TotalQueries     int
	FailedQueries    int
	TotalExecutions  int
	FailedExecutions int
	BytesSent        int
	BytesReceived    int
}

// ConnectionConfig holds configuration options for establishing a database connection.
type ConnectionConfig struct {
	// Host is the database server hostname or IP address.
	Host string

	// Port is the database server port.
	Port int

	// Username for authentication.
	Username string

	// Password for authentication.
	Password string

	// DatabaseName is the name of the database to connect to.
	DatabaseName string

	// TLSConfig holds TLS/SSL configuration.
	TLSConfig *TLSConfig

	// SSHConfig holds SSH tunnel configuration.
	SSHConfig *SSHConfig

	// Timeout is the maximum time to wait for a connection to be established.
	Timeout time.Duration
}

// TLSConfig holds TLS/SSL configuration.
type TLSConfig struct {
	IsEnabled            bool
	CACertPath           string
	ClientCertPath       string
	ClientKeyPath        string
	ClientCertPassword   string
	IsServerNameInCert   bool
	IsInsecureSkipVerify bool
}

// SSHConfig holds SSH tunnel configuration.
type SSHConfig struct {
	IsEnabled bool
	Host      string
	Port      int
	Username  string
	Password  string
}
