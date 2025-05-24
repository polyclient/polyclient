// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package netutil

import (
	"net"
	"strconv"
)

var defaultListen = net.Listen

// IsPortAvailable checks if a port is available on the local machine.
func IsPortAvailable(port int) bool {
	listener, err := defaultListen("tcp", ":"+strconv.Itoa(port))
	if err == nil {
		return true
	}

	if listener != nil {
		_ = listener.Close()
	}

	return false
}

// GetAvailablePort finds an available port on the local machine.
func GetAvailablePort() (int, error) {
	listener, err := defaultListen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)

	return addr.Port, nil
}
