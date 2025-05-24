// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package netutil

import (
	"errors"
	"fmt"
	"net"
)

var defaultIfaceAddrs = func() ([]net.Addr, error) {
	return net.InterfaceAddrs()
}

// LocalIPv4 returns the first non-loopback IPv4 address.
func LocalIPv4() (string, error) {
	addrs, err := defaultIfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("failed to get interface addresses: %w", err)
	}

	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() || ipnet.IP.To4() == nil {
			continue
		}

		return ipnet.IP.String(), nil
	}

	return "", errors.New("no non-loopback IPv4 address found")
}

// LocalIPv6 returns the first non-loopback IPv6 address.
func LocalIPv6() (string, error) {
	addrs, err := defaultIfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("failed to get interface addresses: %w", err)
	}

	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok || ipnet.IP.IsLoopback() {
			continue
		}

		if ipnet.IP.To4() == nil && ipnet.IP.To16() != nil && !ipnet.IP.IsLinkLocalUnicast() {
			return ipnet.IP.String(), nil
		}
	}

	return "", errors.New("no non-loopback IPv6 address found")
}
