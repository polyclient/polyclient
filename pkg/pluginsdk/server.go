package pluginsdk

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	pb "github.com/polyclient/polyclient/proto"
	"google.golang.org/grpc"
)

// pluginServer defines the gRPC server for the plugin.
type pluginServer struct {
	pb.UnimplementedPluginServer
	plugin *Plugin
	mu     sync.RWMutex
}

// Execute executes the specified action handler for a plugin.
func (s *pluginServer) Execute(ctx context.Context, req *pb.PluginExecuteRequest) (res *pb.PluginExecuteResponse, err error) {
	defer func() {
		// We recover from panics in the action handler to prevent the entire
		// gRPC server from crashing. This is important because a panic in the
		// plugin should not affect the main PolyClient process.
		if r := recover(); r != nil {
			err = fmt.Errorf("handler panic: %v", r)
			res = &pb.PluginExecuteResponse{Error: err.Error()}
		}
	}()

	s.mu.RLock()
	handler, exists := s.plugin.Handlers[req.Action]
	s.mu.RUnlock()

	if !exists {
		return &pb.PluginExecuteResponse{Error: "unknown action: " + req.Action}, nil
	}

	result, err := handler(req.Payload)
	if err != nil {
		return &pb.PluginExecuteResponse{Error: err.Error()}, err
	}

	return &pb.PluginExecuteResponse{Payload: result}, nil
}

// Serve starts a gRPC server for the given Plugin. It determines the appropriate communication
// endpoint based on the operating system. On Unix-like systems, it uses a Unix domain socket,
// which is a special file in the file system used for fast, local inter-process communication (IPC).
// On Windows, it uses a named pipe, which is the Windows equivalent for IPC.
//
// Sockets (Unix domain sockets) and pipes (named pipes on Windows) are used here because they
// provide efficient, local communication without the overhead of TCP/IP networking.
func Serve(plugin *Plugin) {
	socketPath := GetSocketPath(plugin.Manifest.Name)

	// Create a listener for the gRPC server
	listener, err := createListener(socketPath)
	if err != nil {
		log.Fatalf("failed to create listener: %v", err)
	}
	defer listener.Close()

	if _, err := os.Stat(socketPath); os.IsNotExist(err) {
		log.Fatalf("socket file was not created: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterPluginServer(grpcServer, &pluginServer{plugin: plugin})

	// Graceful shutdown setup
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	// Server in goroutine so we can block on context
	go func() {
		log.Printf("Server starting on socket: %s", socketPath)

		if err := grpcServer.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()

	// Graceful stop with timeout
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// First try graceful stop
	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Println("Server stopped gracefully")
	case <-shutdownCtx.Done():
		log.Println("Forcing server stop after timeout")
		grpcServer.Stop()
	}
}
