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
	Addr string
	Port int

	httpServer *http.Server
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
	}
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
