// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/polyclient/polyclient/internal/engine"
	"github.com/polyclient/polyclient/internal/netutil"
)

// Server is the HTTP server for the API.
type Server struct {
	Host       string
	Port       int
	Router     *Router
	HTTPServer *http.Server
}

// ServerOption is a function that modifies Server.
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
var defaultServerOptions = func(e *engine.Engine) *Server {
	return &Server{
		Host: e.Settings.API.Host,
		Port: e.Settings.API.Port,
	}
}

// NewServer creates a new HTTP server for the API.
func NewServer(ctx context.Context, e *engine.Engine, options ...ServerOption) (*Server, error) {
	opts := defaultServerOptions(e)
	for _, opt := range options {
		opt(opts)
	}

	var port int

	if netutil.IsPortAvailable(opts.Port) {
		port = opts.Port
	} else {
		foundPort, err := netutil.GetAvailablePort()
		if err != nil {
			return nil, fmt.Errorf("failed to find an available port: %w", err)
		}

		port = foundPort
	}

	router := NewRouter(ctx, e)
	router.RegisterAPIV1Routes()

	if e.Settings.GUI.Enabled {
		router.RegisterGUIRoutes()
	}

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", opts.Host, port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if e.Settings.API.Timeouts.Enabled {
		httpServer.ReadTimeout = time.Duration(e.Settings.API.Timeouts.ReadSeconds)
		httpServer.WriteTimeout = time.Duration(e.Settings.API.Timeouts.WriteSeconds)
		httpServer.IdleTimeout = time.Duration(e.Settings.API.Timeouts.IdleSeconds)
	}

	return &Server{
		Host:       opts.Host,
		Port:       port,
		Router:     router,
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
