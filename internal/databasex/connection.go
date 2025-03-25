// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package databasex

import (
	"context"
	"time"
)

// Querier executes queries that return rows.
type Querier interface {
	Query(ctx context.Context, query string, args ...any) (Rows, error)
}

// Executor runs queries that don't return rows.
type Executor interface {
	Execute(ctx context.Context, query string, args ...any) (Result, error)
}

// TransactionStarter initializes database transactions.
type TransactionStarter interface {
	BeginTransaction(ctx context.Context, opts *TransactionOptions) (Transaction, error)
}

// InfoProvider supplies metadata about the database connection.
type InfoProvider interface {
	ConnectionInfo() ConnectionInfo
}

// Pingable checks the health of a database connection.
type Pingable interface {
	Ping(ctx context.Context) error
}

// Closable closes the database connection.
type Closable interface {
	Close()
}

// RawAccessor exposes the underlying database connection.
type RawAccessor interface {
	Raw() any
}

// Connection aggregates database capabilities into a single interface.
type Connection interface {
	Querier
	Executor
	TransactionStarter
	InfoProvider
	Pingable
	Closable
	RawAccessor
}

// ConnectionInfo holds details about a database connection.
type ConnectionInfo struct {
	DatabaseName  string          `json:"database_name"`
	ServerVersion string          `json:"server_version"`
	ServerType    string          `json:"server_type"`
	Host          string          `json:"host"`
	Port          int             `json:"port"`
	Username      string          `json:"username"`
	IsEncrypted   bool            `json:"is_encrypted"`
	IsSSHEnabled  bool            `json:"is_ssh_enabled"`
	IsConnected   bool            `json:"is_connected"`
	ConnectedAt   time.Time       `json:"connected_at"`
	Stats         ConnectionStats `json:"stats"`
}

// ConnectionStats tracks runtime metrics for a connection.
type ConnectionStats struct {
	Uptime           time.Duration `json:"uptime"`
	TotalQueries     int           `json:"total_queries"`
	FailedQueries    int           `json:"failed_queries"`
	TotalExecutions  int           `json:"total_executions"`
	FailedExecutions int           `json:"failed_executions"`
	BytesSent        int           `json:"bytes_sent"`
	BytesReceived    int           `json:"bytes_received"`
}

// ConnectionConfig defines parameters for establishing a database connection.
type ConnectionConfig struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	Username     string        `json:"username"`
	Password     string        `json:"password"`
	DatabaseName string        `json:"database_name"`
	TLSConfig    *TLSConfig    `json:"tls_config,omitempty"`
	SSHConfig    *SSHConfig    `json:"ssh_config,omitempty"`
	Timeout      time.Duration `json:"timeout"`
}

// TLSConfig manages TLS/SSL settings.
type TLSConfig struct {
	Enabled            bool   `json:"enabled"`
	CACertPath         string `json:"ca_cert_path"`
	ClientCertPath     string `json:"client_cert_path"`
	ClientKeyPath      string `json:"client_key_path"`
	ClientCertPassword string `json:"client_cert_password"`
	VerifyServerName   bool   `json:"verify_server_name"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify"`
}

// SSHConfig defines SSH tunnel settings.
type SSHConfig struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}
