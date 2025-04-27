package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/polyclient/polyclient/internal/constant"
	"github.com/polyclient/polyclient/internal/version"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Format represents the output format.
type Format string

const (
	// JSONFormat makes the logger output JSON.
	JSONFormat Format = "json"

	// TextFormat makes the logger output plain text.
	TextFormat Format = "text"
)

// String returns the string representation of the format.
func (f Format) String() string {
	return string(f)
}

// Config represents the configuration for the logger.
type Config struct {
	Format   Format
	Level    slog.Level
	Rotation RotationConfig
}

// RotationConfig represents the configuration for log rotation.
type RotationConfig struct {
	Enabled    bool
	Filename   string
	MaxSizeMB  int
	MaxAgeDays int
	MaxBackups int
	LocalTime  bool
	Compress   bool
}

// defaultConfig returns the default logger configuration.
var defaultConfig = func() *Config {
	return &Config{
		Level:    slog.LevelInfo,
		Format:   JSONFormat,
		Rotation: *defaultRotationConfig(),
	}
}

// defaultRotationConfig returns the default rotation configuration.
var defaultRotationConfig = func() *RotationConfig {
	return &RotationConfig{
		Enabled:    true,
		Filename:   getDefaultLogsPath(),
		MaxSizeMB:  100,   // 100 MB
		MaxAgeDays: 30,    // keep 30 days
		MaxBackups: 5,     // keep 5 backups
		LocalTime:  true,  // use local time instead of UTC
		Compress:   false, // compress rotated logs
	}
}

// Option is a function that configures the logger.
type Option func(*Config)

// WithFormat sets the output format for the logger.
func WithFormat(format Format) Option {
	return func(cfg *Config) {
		cfg.Format = format
	}
}

// WithLevel sets the log level for the logger.
func WithLevel(level slog.Level) Option {
	return func(cfg *Config) {
		cfg.Level = level
	}
}

// WithRotationEnabled enables or disables log rotation.
func WithRotationEnabled(enabled bool) Option {
	return func(cfg *Config) {
		cfg.Rotation.Enabled = enabled
	}
}

// WithRotationFilename sets the filename for rotated log files.
// Defaults to a platform-specific path (see getDefaultLogsPath).
func WithRotationFilename(filename string) Option {
	return func(cfg *Config) {
		if filename != "" {
			cfg.Rotation.Filename = filename
		}
	}
}

// WithRotationMaxSize sets the maximum size in megabytes of the log file before it gets rotated.
func WithRotationMaxSize(megabytes int) Option {
	return func(cfg *Config) {
		if megabytes > 0 {
			cfg.Rotation.MaxSizeMB = megabytes
		}
	}
}

// WithRotationMaxAge sets the maximum number of days to retain old log files.
func WithRotationMaxAge(days int) Option {
	return func(cfg *Config) {
		if days > 0 {
			cfg.Rotation.MaxAgeDays = days
		}
	}
}

// WithRotationMaxBackups sets the maximum number of old log files to retain.
func WithRotationMaxBackups(count int) Option {
	return func(cfg *Config) {
		if count > 0 {
			cfg.Rotation.MaxBackups = count
		}
	}
}

// WithRotationLocalTime determines if rotated log files should use local time instead of UTC.
func WithRotationLocalTime(useLocal bool) Option {
	return func(cfg *Config) {
		cfg.Rotation.LocalTime = useLocal
	}
}

// WithRotationCompress determines if rotated log files should be compressed using gzip.
func WithRotationCompress(compress bool) Option {
	return func(cfg *Config) {
		cfg.Rotation.Compress = compress
	}
}

// Logger is a custom logger implementation on top of slog.
type Logger struct {
	*slog.Logger
	writers []io.Writer
	closer  io.Closer
}

// NewLogger creates a new Logger with the given options.
func NewLogger(opts ...Option) (*Logger, error) {
	cfg := defaultConfig()

	for _, opt := range opts {
		opt(cfg)
	}

	if err := os.MkdirAll(filepath.Dir(cfg.Rotation.Filename), 0o750); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	writers := make([]io.Writer, 0, 2)
	if version.Version() != "dev" {
		writers = append(writers, os.Stdout)
	}

	var closer io.Closer

	if cfg.Rotation.Enabled {
		rotator := &lumberjack.Logger{
			Filename:   cfg.Rotation.Filename,
			MaxSize:    cfg.Rotation.MaxSizeMB,
			MaxAge:     cfg.Rotation.MaxAgeDays,
			MaxBackups: cfg.Rotation.MaxBackups,
			LocalTime:  cfg.Rotation.LocalTime,
			Compress:   cfg.Rotation.Compress,
		}

		writers = append(writers, rotator)

		closer = rotator
	}

	multiwriter := io.MultiWriter(writers...)

	var handler slog.Handler
	handlerOpts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: !version.IsProd(),
	}

	switch cfg.Format {
	case JSONFormat:
		handler = slog.NewJSONHandler(multiwriter, handlerOpts)
	case TextFormat:
		handler = slog.NewTextHandler(multiwriter, handlerOpts)
	default:
		handler = slog.NewJSONHandler(multiwriter, handlerOpts)
	}

	logger := &Logger{
		Logger:  slog.New(handler),
		writers: writers,
		closer:  closer,
	}

	logger.Info("Logger initialized",
		"app_version", version.Version(),
		"log_path", cfg.Rotation.Filename,
		"level", cfg.Level.String(),
		"format", cfg.Format.String(),
	)

	return logger, nil
}

// Close closes the logger resources and flushes any pending log entries.
func (l *Logger) Close() error {
	return l.closer.Close()
}

// getDefaultLogsPath returns the default path for log files based on the environment and OS.
func getDefaultLogsPath() string {
	if !version.IsProd() {
		tmpDir := filepath.Join("tmp", "logs")
		return filepath.Join(tmpDir, constant.AppName+".log")
	}

	// Try XDG_STATE_HOME first because it's the latest spec
	// See https://specifications.freedesktop.org/basedir-spec/latest/#variables
	if xdgState := os.Getenv("XDG_STATE_HOME"); xdgState != "" {
		return filepath.Join(xdgState, constant.AppName, "logs", constant.AppName+".log")
	}

	switch runtime.GOOS {
	case "windows":
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			return filepath.Join(localAppData, constant.AppName, "logs", constant.AppName+".log")
		}

	case "darwin":
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, "Library", "Logs", constant.AppName, constant.AppName+".log")
		}

	default: // Linux and others
		if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData != "" {
			return filepath.Join(xdgData, constant.AppName, "logs", constant.AppName+".log")
		}

		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, ".local", "share", constant.AppName, "logs", constant.AppName+".log")
		}
	}

	// Ultimate fallback
	return filepath.Join(os.TempDir(), "logs", constant.AppName+".log")
}
