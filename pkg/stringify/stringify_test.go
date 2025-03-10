// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package stringify_test

import (
	"errors"
	"html"
	"testing"
	"time"

	"github.com/polyclient/polyclient/pkg/stringify"
	"github.com/polyclient/polyclient/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestStringify(t *testing.T) {
	mockTime := time.Date(2025, 3, 9, 14, 30, 0, 0, time.UTC)

	t.Parallel()

	t.Run("Handles nil values", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "", stringify.Stringify(nil))
	})

	t.Run("Handles primitive types", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "hello", stringify.Stringify("hello"))
		assert.Equal(t, "42", stringify.Stringify(42))
		assert.Equal(t, "3.14", stringify.Stringify(3.14))
		assert.Equal(t, "true", stringify.Stringify(true))
		assert.Equal(t, "", stringify.Stringify(nil))
	})

	t.Run("Handles time.Time", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "2025-03-09T14:30:00Z", stringify.Stringify(mockTime, stringify.WithDateFormat(time.RFC3339)))
		assert.Equal(t, "09 Mar 25 14:30 UTC", stringify.Stringify(mockTime, stringify.WithDateFormat(time.RFC822)))
		assert.Equal(t, "2025-03-09", stringify.Stringify(mockTime, stringify.WithDateFormat("2006-01-02")))
		assert.Equal(t, "14:30:00", stringify.Stringify(mockTime, stringify.WithDateFormat("15:04:05")))
	})

	t.Run("Handles *time.Time", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "2025-03-09T14:30:00Z", stringify.Stringify(&mockTime, stringify.WithDateFormat(time.RFC3339)))
		assert.Equal(t, "09 Mar 25 14:30 UTC", stringify.Stringify(&mockTime, stringify.WithDateFormat(time.RFC822)))
		assert.Equal(t, "2025-03-09", stringify.Stringify(&mockTime, stringify.WithDateFormat("2006-01-02")))
		assert.Equal(t, "14:30:00", stringify.Stringify(&mockTime, stringify.WithDateFormat("15:04:05")))
	})

	t.Run("Handles fmt.Stringer", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "Hello, my name is Jane and I am 30 years old", stringify.Stringify(mocks.StringerPersonMock{Name: "Jane", Age: 30}))
	})

	t.Run("Handles errors", func(t *testing.T) {
		t.Parallel()

		err := errors.New("something went wrong")
		assert.Equal(t, "something went wrong", stringify.Stringify(err))
	})

	t.Run("Handles slices", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "[]", stringify.Stringify([]any{}))
		assert.Equal(t, "[1, 2, three]", stringify.Stringify([]any{1, 2, "three"}))
	})

	t.Run("Handles maps", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "{}", stringify.Stringify(map[string]any{}))
		assert.Equal(t, "{a:1, b:2}", stringify.Stringify(map[string]any{"a": 1, "b": 2}))
	})

	t.Run("Handles custom formatter", func(t *testing.T) {
		t.Parallel()

		greeterFunc := func(v string) string {
			return "hello " + v
		}

		sanitizerFunc := func(v string) string {
			return html.EscapeString(v)
		}

		assert.Equal(t, "hello Jane", stringify.Stringify("Jane", stringify.WithCustomFormatter(greeterFunc)))

		assert.Equal(t, "&lt;script&gt;console.log(&#39;performing evil XSS attack&#39;)&lt;/script&gt;",
			stringify.Stringify("<script>console.log('performing evil XSS attack')</script>",
				stringify.WithCustomFormatter(sanitizerFunc),
			))
	})
}
