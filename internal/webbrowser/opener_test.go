package webbrowser

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenURL(t *testing.T) {
	t.Parallel()

	originalRunCommand := defaultRunCommand
	originalGetOS := defaultGOOS
	originalIsWSL := defaultIsWSL

	t.Cleanup(func() {
		defaultRunCommand = originalRunCommand
		defaultGOOS = originalGetOS
		defaultIsWSL = originalIsWSL
	})

	tests := []struct {
		name        string
		url         string
		mockOS      string
		mockWSL     bool
		mockRunErr  error
		wantErr     bool
		wantErrMsg  string
		wantCommand string
		wantArgs    []string
	}{
		{
			name:       "empty URL",
			url:        "",
			wantErr:    true,
			wantErrMsg: "URL cannot be empty",
		},
		{
			name:       "invalid URL",
			url:        "not a url",
			wantErr:    true,
			wantErrMsg: "invalid URL",
		},
		{
			name:        "Linux with valid URL",
			url:         "https://example.com",
			mockOS:      "linux",
			wantCommand: "xdg-open",
			wantArgs:    []string{"https://example.com"},
		},
		{
			name:        "macOS with valid URL",
			url:         "https://example.com",
			mockOS:      "darwin",
			wantCommand: "open",
			wantArgs:    []string{"https://example.com"},
		},
		{
			name:        "Windows with valid URL",
			url:         "https://example.com",
			mockOS:      "windows",
			wantCommand: "cmd.exe",
			wantArgs:    []string{"/c", "start", "", "https://example.com"},
		},
		{
			name:        "WSL with valid URL",
			url:         "https://example.com",
			mockOS:      "linux",
			mockWSL:     true,
			wantCommand: "powershell.exe",
			wantArgs:    []string{"Start-Process", "https://example.com"},
		},
		{
			name:       "Linux command fails",
			url:        "https://example.com",
			mockOS:     "linux",
			mockRunErr: errors.New("command failed"),
			wantErr:    true,
			wantErrMsg: "command failed",
		},
		{
			name:       "unsupported OS",
			url:        "https://example.com",
			mockOS:     "plan9",
			wantErr:    true,
			wantErrMsg: "unsupported operating system",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			defaultRunCommand = func(ctx context.Context, name string, args ...string) error {
				if tt.wantCommand != "" && name != tt.wantCommand {
					require.Error(t, errors.New("unexpected command"))
				}

				if len(tt.wantArgs) > 0 {
					for i, arg := range tt.wantArgs {
						if i >= len(args) || args[i] != arg {
							require.Error(t, errors.New("unexpected argument"))
						}
					}
				}

				return tt.mockRunErr
			}

			defaultGOOS = func() string {
				return tt.mockOS
			}

			defaultIsWSL = func() bool {
				return tt.mockWSL
			}

			err := OpenURL(context.Background(), tt.url)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
