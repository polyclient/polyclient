// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package netutil

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalIPv4(t *testing.T) {
	t.Parallel()

	originalIfaceAddrs := defaultIfaceAddrs

	t.Cleanup(func() {
		defaultIfaceAddrs = originalIfaceAddrs
	})

	tests := []struct {
		name        string
		mockAddrs   []net.Addr
		mockErr     error
		expectedIP  string
		expectedErr string
	}{
		{
			name: "success with single IPv4",
			mockAddrs: []net.Addr{
				&net.IPNet{IP: net.ParseIP("192.168.1.1")},
			},
			expectedIP: "192.168.1.1",
		},
		{
			name: "success with multiple addresses",
			mockAddrs: []net.Addr{
				&net.IPNet{IP: net.ParseIP("127.0.0.1")}, // loopback
				&net.IPNet{IP: net.ParseIP("192.168.1.1")},
				&net.IPNet{IP: net.ParseIP("10.0.0.1")},    // should pick first non-loopback
				&net.IPNet{IP: net.ParseIP("2001:db8::1")}, // IPv6
			},
			expectedIP: "192.168.1.1",
		},
		{
			name: "no IPv4 addresses",
			mockAddrs: []net.Addr{
				&net.IPNet{IP: net.ParseIP("127.0.0.1")},   // loopback
				&net.IPNet{IP: net.ParseIP("2001:db8::1")}, // IPv6
			},
			expectedErr: "no non-loopback IPv4 address found",
		},
		{
			name:        "interface error",
			mockErr:     errors.New("interface error"),
			expectedErr: "failed to get interface addresses: interface error",
		},
		{
			name: "non-IPNet address type",
			mockAddrs: []net.Addr{
				&net.UnixAddr{}, // invalid type
			},
			expectedErr: "no non-loopback IPv4 address found",
		},
		{
			name: "IPv6 only",
			mockAddrs: []net.Addr{
				&net.IPNet{IP: net.ParseIP("2001:db8::1")},
			},
			expectedErr: "no non-loopback IPv4 address found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			defaultIfaceAddrs = func() ([]net.Addr, error) {
				return tt.mockAddrs, tt.mockErr
			}

			ip, err := LocalIPv4()
			if tt.expectedErr != "" {
				require.ErrorContains(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedIP, ip)
		})
	}
}

func TestLocalIPv6(t *testing.T) {
	t.Parallel()

	originalIfaceAddrs := defaultIfaceAddrs

	t.Cleanup(func() {
		defaultIfaceAddrs = originalIfaceAddrs
	})

	tests := []struct {
		name        string
		mockAddrs   []net.Addr
		mockErr     error
		expectedIP  string
		expectedErr string
	}{
		{
			name: "success with single IPv6",
			mockAddrs: []net.Addr{
				&net.IPNet{IP: net.ParseIP("2001:db8::1")},
			},
			expectedIP: "2001:db8::1",
		},
		{
			name: "success with multiple addresses",
			mockAddrs: []net.Addr{
				&net.IPNet{IP: net.ParseIP("127.0.0.1")},   // loopback
				&net.IPNet{IP: net.ParseIP("192.168.1.1")}, // IPv4
				&net.IPNet{IP: net.ParseIP("2001:db8::1")},
				&net.IPNet{IP: net.ParseIP("fe80::1")}, // link-local
			},
			expectedIP: "2001:db8::1",
		},
		{
			name: "no IPv6 addresses",
			mockAddrs: []net.Addr{
				&net.IPNet{IP: net.ParseIP("127.0.0.1")},   // loopback
				&net.IPNet{IP: net.ParseIP("192.168.1.1")}, // IPv4
			},
			expectedErr: "no non-loopback IPv6 address found",
		},
		{
			name:        "interface error",
			mockErr:     errors.New("interface error"),
			expectedErr: "failed to get interface addresses: interface error",
		},
		{
			name: "non-IPNet address type",
			mockAddrs: []net.Addr{
				&net.UnixAddr{}, // invalid type
			},
			expectedErr: "no non-loopback IPv6 address found",
		},
		{
			name: "only link-local IPv6",
			mockAddrs: []net.Addr{
				&net.IPNet{IP: net.ParseIP("fe80::1")},
			},
			expectedErr: "no non-loopback IPv6 address found",
		},
		{
			name: "IPv4 only",
			mockAddrs: []net.Addr{
				&net.IPNet{IP: net.ParseIP("192.168.1.1")},
			},
			expectedErr: "no non-loopback IPv6 address found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			defaultIfaceAddrs = func() ([]net.Addr, error) {
				return tt.mockAddrs, tt.mockErr
			}

			ip, err := LocalIPv6()
			if tt.expectedErr != "" {
				require.ErrorContains(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedIP, ip)
		})
	}
}
