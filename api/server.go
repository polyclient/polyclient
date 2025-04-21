// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package api

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/polyclient/polyclient/internal/application"
)

// Server is the HTTP server for the API.
type Server struct {
	Addr string
	Port int

	handler http.Handler
}

// ServerOptions is the configuration options for the API server.
type ServerOptions struct {
	Port int
}

// ServerOption is a function that configures the API server.
type ServerOption func(*ServerOptions)

// WithPort sets the port for the API server.
func WithPort(port int) ServerOption {
	return func(opts *ServerOptions) {
		opts.Port = port
	}
}

// defaultOptions returns the default options for the API server.
var defaultOptions = func() *ServerOptions {
	return &ServerOptions{
		Port: 8080,
	}
}

// NewServer creates a new HTTP server for the API.
func NewServer(app *application.Application, options ...ServerOption) (*Server, error) {
	config := defaultOptions()
	for _, opt := range options {
		opt(config)
	}

	router, err := NewRouter(app)
	if err != nil {
		return nil, fmt.Errorf("failed to create router: %w", err)
	}

	var port int
	if isPortTaken(config.Port) {
		foundPort, err := findAvailablePort()
		if err != nil {
			return nil, fmt.Errorf("failed to find an available port: %w", err)
		}
		port = foundPort
	} else {
		port = config.Port
	}

	addr := fmt.Sprintf("127.0.0.1:%d", port)

	return &Server{
		Addr:    addr,
		Port:    port,
		handler: router,
	}, nil
}

// ListenAndServe starts the HTTP server and listens for requests.
func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.Addr, s.handler)
}

// Shutdown gracefully shuts down the HTTP server.
func (s *Server) Shutdown(ctx context.Context) error {
	server := &http.Server{Addr: s.Addr, Handler: s.handler}
	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()
	return server.ListenAndServe()
}

// isPortTaken checks if a port is already in use.
func isPortTaken(port int) bool {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return true
	}
	defer listener.Close()
	return false
}

// findAvailablePort finds an available port.
func findAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}
