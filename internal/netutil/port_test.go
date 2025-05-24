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

type mockListener struct{}

func (m *mockListener) Accept() (net.Conn, error) { return nil, errors.New("something went wrong") }
func (m *mockListener) Close() error              { return nil }
func (m *mockListener) Addr() net.Addr            { return &net.TCPAddr{} }

func TestIsPortAvailable(t *testing.T) {
	t.Parallel()

	originalListen := defaultListen

	t.Cleanup(func() {
		defaultListen = originalListen
	})

	tests := []struct {
		name        string
		port        int
		listenError error
		listener    net.Listener
		expected    bool
	}{
		{
			name:        "port unavailable (error)",
			port:        80,
			listenError: errors.New("port in use"),
			listener:    nil,
			expected:    false,
		},
		{
			name:        "port available",
			port:        8080,
			listenError: nil,
			listener:    &mockListener{},
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			defaultListen = func(network, address string) (net.Listener, error) {
				return tt.listener, tt.listenError
			}

			result := IsPortAvailable(tt.port)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetAvailablePort(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		port, err := GetAvailablePort()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, port, 1024, "should be above well-known ports")
		assert.LessOrEqual(t, port, 65535, "should be below max port number")
	})

	t.Run("listen error", func(t *testing.T) {
		t.Parallel()

		originalListen := defaultListen

		t.Cleanup(func() {
			defaultListen = originalListen
		})

		defaultListen = func(network, address string) (net.Listener, error) {
			return nil, errors.New("listen failed")
		}

		port, err := GetAvailablePort()

		require.Error(t, err)
		assert.Equal(t, 0, port)
	})
}
