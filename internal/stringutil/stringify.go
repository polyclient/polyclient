// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package stringutil

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Config holds configuration options for converting values to strings.
type StringifyConfig struct {
	// DateFormat specifies the format for time values (default: RFC3339).
	DateFormat string
	// CustomFormatter allows specifying a custom conversion function (default: fmt.Sprintf).
	CustomFormatter func(string) string
	// NilValue specifies how nil values are formatted (default: "").
	NilValue string
	// SliceDelimiters specifies the start/end delimiters for slices (default: "[", "]").
	SliceDelimiters [2]string
	// MapDelimiters specifies the start/end delimiters for maps (default: "{", "}").
	MapDelimiters [2]string
}

// DefaultConfig returns the default configuration.
func StringifyDefaultConfig() *StringifyConfig {
	return &StringifyConfig{
		DateFormat:      time.RFC3339,
		CustomFormatter: func(v string) string { return v },
		NilValue:        "",
		SliceDelimiters: [2]string{"[", "]"},
		MapDelimiters:   [2]string{"{", "}"},
	}
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

// WithNilValue sets the string representation for nil values.
func WithNilValue(val string) StringifyOption {
	return func(cfg *StringifyConfig) {
		cfg.NilValue = val
	}
}

// Stringify converts any value to a string, respecting the provided configuration.
func Stringify(v any, opts ...StringifyOption) string {
	cfg := StringifyDefaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.CustomFormatter == nil {
		cfg.CustomFormatter = func(v string) string { return v }
	}

	return stringifyInternal(v, cfg)
}

// stringifyInternal is the internal implementation for string conversion.
func stringifyInternal(v any, cfg *StringifyConfig) string {
	if v == nil {
		return cfg.NilValue
	}

	switch v := v.(type) {
	case time.Time:
		return v.Format(cfg.DateFormat)

	case *time.Time:
		if v == nil {
			return cfg.NilValue
		}

		return v.Format(cfg.DateFormat)

	case fmt.Stringer:
		return v.String()

	case error:
		return v.Error()
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		var b strings.Builder

		b.WriteString(cfg.SliceDelimiters[0])

		for i := range rv.Len() {
			if i > 0 {
				b.WriteString(", ")
			}

			b.WriteString(stringifyInternal(rv.Index(i).Interface(), cfg))
		}

		b.WriteString(cfg.SliceDelimiters[1])

		return b.String()

	case reflect.Map:
		var b strings.Builder

		b.WriteString(cfg.MapDelimiters[0])

		keys := rv.MapKeys()
		for i, key := range keys {
			if i > 0 {
				b.WriteString(", ")
			}

			keyStr := stringifyInternal(key.Interface(), cfg)
			valStr := stringifyInternal(rv.MapIndex(key).Interface(), cfg)

			if _, err := fmt.Fprintf(&b, "%v:%v", keyStr, valStr); err != nil {
				return err.Error()
			}
		}

		b.WriteString(cfg.MapDelimiters[1])

		return b.String()

	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return rv.Type().String()
	}

	return cfg.CustomFormatter(fmt.Sprintf("%v", v))
}
