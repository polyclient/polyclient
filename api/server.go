// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/polyclient/polyclient/internal/application"
)

// Server is the HTTP server for the API.
type Server struct {
	Host       string
	Port       int
	HTTPServer *http.Server
}

// ServerOption configures the API server.
type ServerOption func(*Server)

// WithHost sets the host for the API server.
func WithHost(host string) ServerOption {
	return func(opts *Server) {
		opts.Host = host
	}
}

// WithPort sets the port for the API server.
func WithPort(port int) ServerOption {
	return func(opts *Server) {
		opts.Port = port
	}
}

// defaultOptions returns the default options for the API server.
var defaultOptions = func() *Server {
	return &Server{
		Host: "127.0.0.1",
		Port: 8080,
	}
}

// NewServer creates a new HTTP server for the API.
func NewServer(app *application.Application, options ...ServerOption) (*Server, error) {
	if app == nil {
		return nil, errors.New("application cannot be nil")
	}

	config := defaultOptions()
	for _, opt := range options {
		opt(config)
	}

	router, err := NewRouter(app)
	if err != nil {
		return nil, fmt.Errorf("failed to create router: %w", err)
	}

	var port int
	if !isPortAvailable(config.Port) {
		foundPort, err := findAvailablePort(config.Host)
		if err != nil {
			return nil, fmt.Errorf("failed to find an available port: %w", err)
		}

		port = foundPort
	} else {
		port = config.Port
	}

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", config.Host, port),
		Handler:           router,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       15 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	return &Server{
		Host:       config.Host,
		Port:       port,
		HTTPServer: httpServer,
	}, nil
}

// ListenAndServe starts the HTTP server and listens for requests.
func (s *Server) ListenAndServe() error {
	if s.HTTPServer == nil {
		return errors.New("http server is not initialized")
	}

	slog.Info("Starting server", "addr", s.HTTPServer.Addr)

	return s.HTTPServer.ListenAndServe()
}

// Shutdown gracefully shuts down the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	if s.HTTPServer == nil {
		return errors.New("http server is not initialized")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	slog.Info("Shutting down server", "addr", s.HTTPServer.Addr)

	return s.HTTPServer.Shutdown(ctx)
}

// isPortAvailable checks if a TCP port is available for use.
func isPortAvailable(port int) bool {
	if port < 1 || port > 65535 {
		return false
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}

	_ = listener.Close()

	return true
}

// findAvailablePort finds an available port.
func findAvailablePort(host string) (int, error) {
	addr := ":0"
	if host != "" {
		addr = host + ":0"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return 0, fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	defer listener.Close()

	tcpAddr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return 0, fmt.Errorf("unexpected address type: %T", listener.Addr())
	}

	return tcpAddr.Port, nil
}
