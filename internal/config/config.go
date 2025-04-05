// SPDX-FileCopyrightText: 2025 The PolyClient Authors
//
// SPDX-License-Identifier: GPL-3.0-or-later WITH LicenseRef-PolyClient-Plugin-Exception

package config

// PolyClientConfig represents the PolyClient configuration structure.
type PolyClientConfig struct {
	Schema   string                 `json:"$schema,omitempty"`
	Theme    string                 `json:"theme,omitempty"`
	Language string                 `json:"language,omitempty"`
	Editor   PolyClientConfigEditor `json:"editor"`
}

// PolyClientConfigEditor represents the PolyClient editor configuration structure.
type PolyClientConfigEditor struct {
	VimMode        bool   `json:"vimMode,omitempty"`
	Minimap        bool   `json:"minimap,omitempty"`
	LineNumbers    string `json:"lineNumbers,omitempty"`
	CursorStyle    string `json:"cursorStyle,omitempty"`
	CursorBlinking string `json:"cursorBlinking,omitempty"`
	FontFamily     string `json:"fontFamily,omitempty"`
	FontSize       int    `json:"fontSize,omitempty"`
	FontWeight     string `json:"fontWeight,omitempty"`
	FontLigatures  string `json:"fontLigatures,omitempty"`
}

// PolyClientConfigHistory represents the PolyClient history configuration structure.
type PolyClientConfigHistory struct {
	Queries PolyClientConfigHistoryQueries `json:"queries"`
}

// PolyClientConfigHistoryQueries represents the PolyClient history query configuration structure.
type PolyClientConfigHistoryQueries struct {
	Enabled  bool `json:"enabled,omitempty"`
	MaxItems int  `json:"maxItems,omitempty"`
}

// PolyClientConfigLayout represents the PolyClient layout configuration structure.
type PolyClientConfigLayout struct {
	Explorer      PolyClientConfigLayoutExplorer  `json:"explorer"`
	Assistant     PolyClientConfigLayoutAssistant `json:"assistant"`
	Output        PolyClientConfigLayoutOutput    `json:"output"`
	RememberSizes bool                            `json:"rememberSizes,omitempty"`
	ZenMode       bool                            `json:"zenMode,omitempty"`
}

// PolyClientConfigLayoutExplorer represents the PolyClient explorer layout configuration structure.
type PolyClientConfigLayoutExplorer struct {
	DefaultSize  int    `json:"defaultSize,omitempty"`
	DockPosition string `json:"dockPosition,omitempty"`
}

// PolyClientConfigLayoutAssistant represents the PolyClient assistant layout configuration structure.
type PolyClientConfigLayoutAssistant struct {
	DefaultSize  int    `json:"defaultSize,omitempty"`
	DockPosition string `json:"dockPosition,omitempty"`
}

// PolyClientConfigLayoutOutput represents the PolyClient output layout configuration structure.
type PolyClientConfigLayoutOutput struct {
	DefaultSize  int    `json:"defaultSize,omitempty"`
	DockPosition string `json:"dockPosition,omitempty"`
}

// PolyClientConfigERD represents the PolyClient ERD configuration structure.
type PolyClientConfigERD struct {
	DefaultLayout   string `json:"defaultLayout,omitempty"`
	ShowPrimaryKeys bool   `json:"showPrimaryKeys,omitempty"`
	ShowForeignKeys bool   `json:"showForeignKeys,omitempty"`
	ShowColumnTypes bool   `json:"showColumnTypes,omitempty"`
	ZoomLevel       int    `json:"zoomLevel,omitempty"`
}

// PolyClientConfigBackup represents the PolyClient backup configuration structure.
type PolyClientConfigBackup struct {
	Storage    PolyClientConfigBackupStorage    `json:"storage"`
	Scheduling PolyClientConfigBackupScheduling `json:"scheduling"`
}

// PolyClientConfigBackupStorage represents the PolyClient backup storage configuration structure.
type PolyClientConfigBackupStorage struct {
	DirectoryPath string `json:"directoryPath,omitempty"`
}

// PolyClientConfigBackupScheduling represents the PolyClient backup scheduling configuration structure.
type PolyClientConfigBackupScheduling struct {
	Enabled        bool   `json:"enabled,omitempty"`
	CronExpression string `json:"cronExpression,omitempty"`
	KeepLastN      int    `json:"keepLastN,omitempty"`
}

// PolyClientConfigA11y represents the PolyClient accessibility configuration structure.
type PolyClientConfigA11y struct {
	ScreenReaderEnabled bool `json:"screenReaderEnabled,omitempty"`
}

// PolyClientConfigSystem represents the PolyClient system configuration structure.
type PolyClientConfigSystem struct {
	Performance PolyClientConfigSystemPerformance `json:"performance,omitempty"`
	Updates     PolyClientConfigSystemUpdates     `json:"updates,omitempty"`
}

// PolyClientConfigSystemPerformance represents the PolyClient system performance configuration structure.
type PolyClientConfigSystemPerformance struct {
	HardwareAccelerationEnabled bool `json:"hardwareAccelerationEnabled,omitempty"`
	ThrottleInactiveTabQueries  bool `json:"throttleInactiveTabQueries,omitempty"`
}

// PolyClientConfigSystemUpdates represents the PolyClient system updates configuration structure.
type PolyClientConfigSystemUpdates struct {
	CheckAutomatically    bool   `json:"checkAutomatically,omitempty"`
	DownloadAutomatically bool   `json:"downloadAutomatically,omitempty"`
	Channel               string `json:"channel,omitempty"`
}
