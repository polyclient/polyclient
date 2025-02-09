//go:build linux || darwin

package pluginsdk

import (
	"fmt"
	"net"
	"os"
)

// createListener creates a Unix domain socket listener.
func createListener(socketPath string) (net.Listener, error) {
	if err := os.Remove("/tmp/db_sqlite.sock"); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to clean up socket: %w", err)
	}

	listener, err := net.Listen("unix", "/tmp/db_sqlite.sock")
	if err != nil {
		return nil, fmt.Errorf("failed to listen on unix socket: %w", err)
	}

	return listener, nil
}

// cleanupSocket removes the Unix domain socket.
func cleanupSocket(socketPath string) error {
	return os.Remove(socketPath)
}
