// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/polyclient/polyclient/api/middleware"
)

type Server struct {
	Addr string
	Port int

	httpServer *http.Server
}

func NewServer() (*Server, error) {
	router := NewRouter()

	port, err := findAvailablePort()
	if err != nil {
		return nil, fmt.Errorf("failed to find an available port: %w", err)
	}

	stack := middleware.CreateStack(
		middleware.Logger,
		middleware.Recover,
	)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf("localhost:%d", port),
		Handler:           stack(router),
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      5 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	return &Server{
		Addr:       httpServer.Addr,
		Port:       port,
		httpServer: httpServer,
	}, nil
}

func (s *Server) ListenAndServe() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func findAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port, nil
}
