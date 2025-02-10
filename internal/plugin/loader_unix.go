//go:build !windows

package plugin

import (
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// createSocketClient creates a gRPC client connection for a Unix domain socket.
func createSocketClientConnection(socketPath string) (*grpc.ClientConn, error) {
	if _, err := os.Stat(socketPath); os.IsNotExist(err) {
		return nil, err
	}

	return grpc.NewClient("unix://"+socketPath,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}

// cleanupSocketClientConnection remopves a Unix domain socket file.
func cleanupSocketClientConnection(socketPath string) error {
	return os.Remove(socketPath)
}
