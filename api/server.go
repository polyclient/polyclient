package api

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/polyclient/polyclient/api/middleware"
)

type Server struct {
	HttpServer *http.Server
}

func NewServer() *Server {
	router := NewRouter()
	port, err := findAvailablePort()
	if err != nil {
		log.Fatalf("Failed to find an available port: %v", err)
	}

	stack := middleware.CreateStack(
		middleware.Logger,
	)

	httpServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           stack(router),
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	return &Server{HttpServer: httpServer}
}

func (s *Server) ListenAndServe() error {
	if err := s.HttpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.HttpServer.Shutdown(ctx)
}

func findAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}
