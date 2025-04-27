// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package settings

import "github.com/go-playground/validator/v10"

// Settings holds the settings for PolyClient.
type Settings struct {
	Drivers Drivers `json:"drivers" validate:"required,dive"`
	API     API     `json:"api" validate:"required,dive"`
	GUI     GUI     `json:"gui" validate:"required,dive"`
}

// Option is a function that configures the settings.
type Option func(*Settings)

func defaultSettings() *Settings {
	return &Settings{
		Drivers: *NewDrivers(),
		API:     *NewAPI(),
		GUI:     *NewGUI(),
	}
}

// NewSettings returns a new Settings instance with the default values.
func NewSettings(opts ...Option) *Settings {
	settings := defaultSettings()

	for _, opt := range opts {
		opt(settings)
	}

	return settings
}

// Validate validates the settings.
func (settings *Settings) Validate() error {
	return validator.New().Struct(settings)
}
