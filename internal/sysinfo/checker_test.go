// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package sysinfo

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsWSL(t *testing.T) {
	t.Parallel()

	originalFileReader := defaultFileReader
	originalGOOS := defaultGOOS

	t.Cleanup(func() {
		defaultFileReader = originalFileReader
		defaultGOOS = originalGOOS
	})

	tests := []struct {
		name          string
		mockOS        string
		mockFile      string
		mockData      []byte
		mockErr       error
		expectedIsWSL bool
	}{
		{
			name:          "WSL2 via /proc/version (Microsoft)",
			mockOS:        "linux",
			mockFile:      "/proc/version",
			mockData:      []byte("Linux version 5.10.60.1-microsoft-standard-WSL2"),
			expectedIsWSL: true,
		},
		{
			name:          "WSL1 via /proc/version (WSL)",
			mockOS:        "linux",
			mockFile:      "/proc/version",
			mockData:      []byte("Linux version 4.4.0-18362-Microsoft (WSL)"),
			expectedIsWSL: true,
		},
		{
			name:          "WSL via /proc/sys/kernel/osrelease",
			mockOS:        "linux",
			mockFile:      "/proc/sys/kernel/osrelease",
			mockData:      []byte("5.15.90.1-microsoft-standard-WSL2"),
			expectedIsWSL: true,
		},
		{
			name:          "Native Linux (No WSL indicators)",
			mockOS:        "linux",
			mockFile:      "/proc/version",
			mockData:      []byte("Linux version 5.15.0-76-generic (Ubuntu)"),
			expectedIsWSL: false,
		},
		{
			name:          "Empty /proc/version file",
			mockOS:        "linux",
			mockFile:      "/proc/version",
			mockData:      []byte(""),
			expectedIsWSL: false,
		},
		{
			name:          "File read error for /proc/version",
			mockOS:        "linux",
			mockFile:      "/proc/version",
			mockErr:       errors.New("permission denied"),
			expectedIsWSL: false,
		},
		{
			name:          "Non-existent file",
			mockOS:        "linux",
			mockFile:      "/proc/version",
			mockErr:       os.ErrNotExist,
			expectedIsWSL: false,
		},
		{
			name:          "Case insensitive WSL detection",
			mockOS:        "linux",
			mockFile:      "/proc/version",
			mockData:      []byte("Linux version 5.10.60.1-MICROSOFT-STANDARD-WSL2"),
			expectedIsWSL: true,
		},
		{
			name:          "WSL with partial file content",
			mockOS:        "linux",
			mockFile:      "/proc/version",
			mockData:      []byte("microsoft"),
			expectedIsWSL: true,
		},
		{
			name:          "Non-Linux OS",
			mockOS:        "windows",
			expectedIsWSL: false,
		},
		{
			name:          "Multiple files, one indicates WSL",
			mockOS:        "linux",
			mockFile:      "/proc/sys/kernel/osrelease",
			mockData:      []byte("5.15.90.1-microsoft-standard-WSL2"),
			expectedIsWSL: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			defaultFileReader = func(name string) ([]byte, error) {
				return tt.mockData, tt.mockErr
			}

			defaultGOOS = func() string {
				return tt.mockOS
			}

			assert.Equal(t, tt.expectedIsWSL, IsWSL())
		})
	}
}
