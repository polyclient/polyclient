//go:build windows

package pluginsdk

import (
	"fmt"
	"net"

	"github.com/Microsoft/go-winio"
)

// createListener creates a named pipe listener on Windows
func createListener(socketPath string) (net.Listener, error) {
	listener, err := winio.ListenPipe(socketPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on named pipe: %w", err)
	}

	return listener, nil
}

// cleanupSocket is a no-op on Windows
func cleanupSocket(socketPath string) error {
	return nil
}
