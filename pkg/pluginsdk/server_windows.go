//go:build windows

package pluginsdk

import (
	"fmt"
	"net"

	"github.com/Microsoft/go-winio"
)

// createListener creates and returns a named pipe listener for the given socket path.
func createListener(socketPath string) (net.Listener, error) {
	listener, err := winio.ListenPipe(socketPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on named pipe: %w", err)
	}

	return listener, nil
}
