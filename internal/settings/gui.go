// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package settings

// GUI holds the settings for the GUI.
type GUI struct {
	Enabled bool   `json:"enabled" validate:"boolean"`
	Path    string `json:"path" validate:"required,filepath"`
	Theme   string `json:"theme" validate:"required,oneof=light dark system"` // TODO: add support for custom themes
	Locale  string `json:"locale" validate:"required,bcp47_language_tag"`
}

// GUIOption is a function that configures the settings for the GUI.
type GUIOption func(*GUI)

// WithGUIEnabled sets the enabled configuration for the GUI.
func WithGUIEnabled(enabled bool) GUIOption {
	return func(gui *GUI) {
		gui.Enabled = enabled
	}
}

// WithGUIPath sets the path configuration for the GUI.
func WithGUIPath(path string) GUIOption {
	return func(gui *GUI) {
		gui.Path = path
	}
}

// WithGUITheme sets the theme configuration for the GUI.
func WithGUITheme(theme string) GUIOption {
	return func(gui *GUI) {
		gui.Theme = theme
	}
}

// WithGUILocale sets the locale configuration for the GUI.
func WithGUILocale(locale string) GUIOption {
	return func(gui *GUI) {
		gui.Locale = locale
	}
}

// defaultGUISettings returns the default settings for the GUI.
func defaultGUISettings() *GUI {
	return &GUI{
		Enabled: true,
		Path:    "/*",
		Theme:   "system",
		Locale:  "en",
	}
}

// NewGUI returns a new GUI instance with the default values.
func NewGUI(opts ...GUIOption) *GUI {
	settings := defaultGUISettings()

	for _, opt := range opts {
		opt(settings)
	}

	return settings
}
