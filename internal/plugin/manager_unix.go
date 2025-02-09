//go:build linux || darwin

package plugin

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// createSocketClient creates a gRPC client connection for a Unix domain socket.
func createSocketClient(socketPath string) (*grpc.ClientConn, error) {
	return grpc.NewClient("unix:///tmp/db_sqlite.sock",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
