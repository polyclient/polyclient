// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package polling_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/polyclient/polyclient/internal/polling"
	"github.com/stretchr/testify/require"
)

func TestUntilReady(t *testing.T) {
	tests := []struct {
		name        string
		handler     http.HandlerFunc
		wantErr     bool
		contextFunc func() (context.Context, context.CancelFunc)
	}{
		{
			name: "returns ok immediately",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			wantErr: false,
			contextFunc: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 2*time.Second)
			},
		},
		{
			name: "retries and times out after short duration",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr: true,
			contextFunc: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 500*time.Millisecond)
			},
		},
		{
			name: "returns early if context is cancelled",
			handler: func(w http.ResponseWriter, r *http.Request) {
				// This handler deliberately sleeps to simulate a slow server
				time.Sleep(500 * time.Millisecond)
				w.WriteHeader(http.StatusOK)
			},
			wantErr: true, // Should get context canceled error
			contextFunc: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
				return ctx, cancel
			},
		},
		{
			name:    "returns ok after some retries",
			handler: newEventuallySuccessHandler(3),
			wantErr: false,
			contextFunc: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 1*time.Second)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			ctx, cancel := tt.contextFunc()
			defer cancel()

			err := polling.UntilReady(ctx, server.URL)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func newEventuallySuccessHandler(successAfterAttempts int) http.HandlerFunc {
	attempts := 0

	return func(w http.ResponseWriter, r *http.Request) {
		attempts++

		if attempts >= successAfterAttempts {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
