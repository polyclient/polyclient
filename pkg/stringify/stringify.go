// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package stringify

import (
	"fmt"
	"reflect"
	"time"

	"github.com/samber/lo"
)

// StringifyConfig holds configuration options for converting values to strings.
type StringifyConfig struct {
	// DateFormat specifies the format for time values (default: RFC3339).
	DateFormat string
	// CustomFormatter allows specifying a custom conversion function (default: fmt.Sprintf).
	CustomFormatter func(string) string
}

// StringifyOption modifies the behavior of the Stringify function.
type StringifyOption func(*StringifyConfig)

// WithDateFormat specifies the format for time values.
func WithDateFormat(format string) StringifyOption {
	return func(cfg *StringifyConfig) {
		cfg.DateFormat = format
	}
}

// WithCustomFormatter allows specifying a custom conversion function.
func WithCustomFormatter(f func(string) string) StringifyOption {
	return func(cfg *StringifyConfig) {
		cfg.CustomFormatter = f
	}
}

// Stringify converts any value to a string, respecting the provided configuration.
func Stringify(v any, opts ...StringifyOption) string {
	cfg := &StringifyConfig{
		DateFormat: time.RFC3339,
		CustomFormatter: func(v string) string {
			return v
		},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	switch v := v.(type) {
	case nil:
		return ""
	case time.Time:
		return v.Format(cfg.DateFormat)
	case *time.Time:
		if v == nil {
			return ""
		}

		return v.Format(cfg.DateFormat)
	case fmt.Stringer:
		return v.String()
	case error:
		return v.Error()
	default:
		// Handle slices, arrays, and maps gracefully
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			var result string

			for i := range rv.Len() {
				if i > 0 {
					result += ", "
				}

				result += Stringify(rv.Index(i).Interface(), opts...)
			}

			return "[" + result + "]"
		case reflect.Map:
			keys := lo.Uniq(lo.Keys(v.(map[string]any)))

			var result string

			for i, key := range keys {
				if i > 0 {
					result += ", "
				}

				result += fmt.Sprintf("%v:%v", key, Stringify(rv.MapIndex(reflect.ValueOf(key)).Interface(), opts...))
			}

			return "{" + result + "}"
		}
	}

	return cfg.CustomFormatter(fmt.Sprintf("%v", v))
}
