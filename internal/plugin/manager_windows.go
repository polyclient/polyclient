//go:build windows

package plugin

import (
	"context"
	"net"

	"github.com/Microsoft/go-winio"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// createSocketClient creates a gRPC client connection for a named pipe on Windows
func createSocketClient(socketPath string) (*grpc.ClientConn, error) {
	return grpc.NewClient(
		"pipe:"+socketPath,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
			return winio.DialPipe(addr, nil)
		}),
	)
}
