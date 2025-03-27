package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Server struct {
	router *Router
	server *http.Server
}

func NewServer() *Server {
	router := NewRouter()

	return &Server{
		router: router,
	}
}

func (s *Server) Run(ctx context.Context) error {
	port, err := findAvailablePort()
	if err != nil {
		return fmt.Errorf("failed to find available port: %w", err)
	}

	s.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           s.router.echo,
		ReadHeaderTimeout: 3 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	fmt.Printf("Starting server on port %d\n", port)

	errChan := make(chan error, 1)
	go func() {
		if err := s.router.echo.StartServer(s.server); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}

		return nil

	case err := <-errChan:
		return err
	}
}

func findAvailablePort() (int, error) {
	port := 8080
	for {
		listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			port++
			continue
		}
		listener.Close()
		return port, nil
	}
}
