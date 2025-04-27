// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package system

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func overrideInterfaceAddrs(t *testing.T, fake func() ([]net.Addr, error)) {
	t.Helper()

	original := interfaceAddrs
	interfaceAddrs = fake

	t.Cleanup(func() {
		interfaceAddrs = original
	})
}

func TestGetLocalIP_Success(t *testing.T) {
	t.Parallel()

	overrideInterfaceAddrs(t, func() ([]net.Addr, error) {
		return []net.Addr{
			&net.IPNet{IP: net.ParseIP("192.168.1.50")},
			&net.IPNet{IP: net.ParseIP("127.0.0.1")}, // loopback, ignored
		}, nil
	})

	ip, err := GetLocalIP()

	require.NoError(t, err)
	require.Equal(t, "192.168.1.50", ip)
}

func TestGetLocalIP_NoNonLoopbackIP(t *testing.T) {
	t.Parallel()

	overrideInterfaceAddrs(t, func() ([]net.Addr, error) {
		return []net.Addr{
			&net.IPNet{IP: net.ParseIP("127.0.0.1")}, // only loopback
		}, nil
	})

	ip, err := GetLocalIP()

	require.Error(t, err)
	require.Empty(t, ip)
}

func TestGetLocalIP_ErrorGettingAddresses(t *testing.T) {
	t.Parallel()

	overrideInterfaceAddrs(t, func() ([]net.Addr, error) {
		return nil, errors.New("error getting addresses")
	})

	ip, err := GetLocalIP()

	require.Error(t, err)
	require.Empty(t, ip)
}
