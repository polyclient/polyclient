// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package stringify_test

import (
	"errors"
	"html"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/polyclient/polyclient/pkg/stringify"
	"github.com/polyclient/polyclient/test/mocks"
)

func TestStringify(t *testing.T) {
	mockTime := time.Date(2025, 3, 9, 14, 30, 0, 0, time.UTC)

	t.Parallel()

	tests := []struct {
		name     string
		input    any
		opts     []stringify.StringifyOption
		expected string
	}{
		{
			name:     "nil value",
			input:    nil,
			expected: "",
		},
		{
			name:     "string",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "integer",
			input:    42,
			expected: "42",
		},
		{
			name:     "float",
			input:    3.14,
			expected: "3.14",
		},
		{
			name:     "boolean",
			input:    true,
			expected: "true",
		},
		{
			name:     "time.Time with RFC3339",
			input:    mockTime,
			opts:     []stringify.StringifyOption{stringify.WithDateFormat(time.RFC3339)},
			expected: "2025-03-09T14:30:00Z",
		},
		{
			name:     "time.Time with RFC822",
			input:    mockTime,
			opts:     []stringify.StringifyOption{stringify.WithDateFormat(time.RFC822)},
			expected: "09 Mar 25 14:30 UTC",
		},
		{
			name:     "*time.Time with RFC3339",
			input:    &mockTime,
			opts:     []stringify.StringifyOption{stringify.WithDateFormat(time.RFC3339)},
			expected: "2025-03-09T14:30:00Z",
		},
		{
			name:     "fmt.Stringer",
			input:    mocks.StringerPersonMock{Name: "Jane", Age: 30},
			expected: "Hello, my name is Jane and I am 30 years old",
		},
		{
			name:     "error",
			input:    errors.New("something went wrong"),
			expected: "something went wrong",
		},
		{
			name:     "empty slice",
			input:    []any{},
			expected: "[]",
		},
		{
			name:     "mixed slice",
			input:    []any{1, 2, "three"},
			expected: "[1, 2, three]",
		},
		{
			name:     "empty map",
			input:    map[string]any{},
			expected: "{}",
		},
		{
			name:     "string map",
			input:    map[string]any{"a": 1},
			expected: "{a:1}",
		},
		{
			name:     "custom nil value",
			input:    nil,
			opts:     []stringify.StringifyOption{stringify.WithNilValue("N/A")},
			expected: "N/A",
		},
		{
			name:  "custom formatter - greeter",
			input: "Jane",
			opts: []stringify.StringifyOption{
				stringify.WithCustomFormatter(func(v string) string { return "hello " + v }),
			},
			expected: "hello Jane",
		},
		{
			name:  "custom formatter - HTML escape",
			input: "<script>alert('xss')</script>",
			opts: []stringify.StringifyOption{
				stringify.WithCustomFormatter(html.EscapeString),
			},
			expected: "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := stringify.Stringify(tt.input, tt.opts...)
			assert.Equal(t, tt.expected, result)
		})
	}
}
