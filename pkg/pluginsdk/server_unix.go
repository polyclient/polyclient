//go:build !windows

package pluginsdk

import (
	"fmt"
	"net"
	"os"
)

// createListener creates and returns a Unix domain socket listener for the given socket path.
func createListener(socketPath string) (net.Listener, error) {
	// Remove old socket if it exists to prevent "address already in use" errors
	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to clean up socket: %w", err)
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on unix socket: %w", err)
	}

	return listener, nil
}
