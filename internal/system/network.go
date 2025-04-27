// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package system

import (
	"errors"
	"net"
)

var interfaceAddrs = net.InterfaceAddrs // internal override for tests

// GetLocalIP returns the first non-loopback IPv4 address of the machine.
func GetLocalIP() (string, error) {
	addrs, err := interfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			return ipnet.IP.String(), nil
		}
	}

	return "", errors.New("no non-loopback IPv4 address found")
}
