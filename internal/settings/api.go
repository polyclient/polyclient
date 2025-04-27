// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package settings

import (
	"time"

	"github.com/polyclient/polyclient/internal/constant"
)

// API holds the settings for the API.
type API struct {
	Host        string         `json:"host" validate:"required,hostname|ip"`
	Port        int            `json:"port" validate:"required,gte=1024,lte=65535"`
	CORS        APICORS        `json:"cors" validate:"required"`
	RateLimit   APIRateLimit   `json:"rateLimit" validate:"required"`
	Compression APICompression `json:"compression" validate:"required"`
	Cache       APICache       `json:"cache" validate:"required"`
	Timeouts    APITimeouts    `json:"timeouts" validate:"required"`
}

// APICORS configures Cross-Origin Resource Sharing.
type APICORS struct {
	Enabled        bool     `json:"enabled" validate:"boolean"`
	AllowedOrigins []string `json:"allowedOrigins" validate:"required,dive,http_url"`
	AllowedMethods []string `json:"allowedMethods" validate:"required,dive,eq=GET|eq=POST|eq=PUT|eq=PATCH|eq=DELETE|eq=OPTIONS"`
	AllowedHeaders []string `json:"allowedHeaders" validate:"required"`
	MaxAge         int      `json:"maxAge" validate:"required,gte=0"`
}

// APIRateLimit configures rate limiting.
type APIRateLimit struct {
	Enabled           bool          `json:"enabled" validate:"boolean"`
	RequestsPerMinute int           `json:"requestsPerMinute" validate:"required,gte=1"`
	WindowLength      time.Duration `json:"windowLength" validate:"required,gte=1m"`
}

// APICompression configures compression.
type APICompression struct {
	Enabled bool `json:"enabled" validate:"boolean"`
	Level   int  `json:"level" validate:"required,gte=0,lte=9"`
}

// APICache configures caching.
type APICache struct {
	Enabled           bool    `json:"enabled" validate:"boolean"`
	DefaultTTLSeconds int     `json:"defaultTTLSeconds" validate:"required,gte=0"`
	MaxSizeMB         float64 `json:"maxSizeMB" validate:"required,gte=0"`
}

// APITimeouts configures request timeouts.
type APITimeouts struct {
	Enabled      bool `json:"enabled" validate:"boolean"`
	ReadSeconds  int  `json:"readSeconds" validate:"required,gte=1"`
	WriteSeconds int  `json:"writeSeconds" validate:"required,gte=1"`
	IdleSeconds  int  `json:"idleSeconds" validate:"required,gte=1"`
}

// APIOption is a function that configures the settings for the API.
type APIOption func(*API)

// WithAPIHost sets the host for the API server.
func WithAPIHost(host string) APIOption {
	return func(api *API) {
		api.Host = host
	}
}

// WithAPIPort sets the port for the API server.
func WithAPIPort(port int) APIOption {
	return func(api *API) {
		api.Port = port
	}
}

// WithAPICORS sets the CORS configuration for the API server.
func WithAPICORS(enabled bool, allowedOrigins, allowedMethods, allowedHeaders []string, maxAge int) APIOption {
	return func(api *API) {
		api.CORS.Enabled = enabled
		api.CORS.AllowedOrigins = allowedOrigins
		api.CORS.AllowedMethods = allowedMethods
		api.CORS.AllowedHeaders = allowedHeaders
		api.CORS.MaxAge = maxAge
	}
}

// WithAPIRateLimit sets the rate limit configuration for the API server.
func WithAPIRateLimit(enabled bool, requestsPerMinute int, windowLength time.Duration) APIOption {
	return func(api *API) {
		api.RateLimit.Enabled = enabled
		api.RateLimit.RequestsPerMinute = requestsPerMinute
		api.RateLimit.WindowLength = windowLength
	}
}

// WithAPICompression sets the compression configuration for the API server.
func WithAPICompression(enabled bool, level int) APIOption {
	return func(api *API) {
		api.Compression.Enabled = enabled
		api.Compression.Level = level
	}
}

// WithAPICache sets the cache configuration for the API server.
func WithAPICache(enabled bool, defaultTTLSeconds int, maxSizeMB float64) APIOption {
	return func(api *API) {
		api.Cache.Enabled = enabled
		api.Cache.DefaultTTLSeconds = defaultTTLSeconds
		api.Cache.MaxSizeMB = maxSizeMB
	}
}

// WithAPITimeouts sets the timeouts configuration for the API server.
func WithAPITimeouts(enabled bool, readSeconds, writeSeconds, idleSeconds int) APIOption {
	return func(api *API) {
		api.Timeouts.Enabled = enabled
		api.Timeouts.ReadSeconds = readSeconds
		api.Timeouts.WriteSeconds = writeSeconds
		api.Timeouts.IdleSeconds = idleSeconds
	}
}

func defaultAPISettings() *API {
	return &API{
		Host: "http://127.0.0.1",
		Port: 8080,
		CORS: APICORS{
			Enabled:        false,
			AllowedOrigins: []string{"http://localhost:*", "http://127.0.0.1:*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Content-Type", constant.HTTPHeaderConnectionName},
			MaxAge:         86400,
		},
		RateLimit: APIRateLimit{
			Enabled:           false,
			RequestsPerMinute: 100,
			WindowLength:      1 * time.Minute,
		},
		Compression: APICompression{
			Enabled: true,
			Level:   6,
		},
		Cache: APICache{
			Enabled:           true,
			DefaultTTLSeconds: 300,
			MaxSizeMB:         50,
		},
		Timeouts: APITimeouts{
			Enabled:      true,
			ReadSeconds:  30,
			WriteSeconds: 30,
			IdleSeconds:  60,
		},
	}
}

// NewAPI returns a new API instance with the default values.
func NewAPI(opts ...APIOption) *API {
	settings := defaultAPISettings()

	for _, opt := range opts {
		opt(settings)
	}

	return settings
}
