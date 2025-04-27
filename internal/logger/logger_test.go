package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/polyclient/polyclient/internal/constant"
	"github.com/polyclient/polyclient/internal/version"
	"github.com/stretchr/testify/suite"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerTestSuite struct {
	suite.Suite
	tempDir string
}

func (s *LoggerTestSuite) SetupSuite() {
	var err error
	s.tempDir, err = os.MkdirTemp("", "logger-test")
	s.Require().NoError(err)
}

func (s *LoggerTestSuite) TearDownSuite() {
	_ = os.RemoveAll(s.tempDir)
}

func (s *LoggerTestSuite) SetupTest() {
	version.SetVersion("test")
}

func TestLoggerTestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

func (s *LoggerTestSuite) TestNewLoggerDefaultConfig() {
	logger, err := NewLogger()
	s.Require().NoError(err)
	s.NotNil(logger)

	s.NotNil(logger.writers)
	s.NotNil(logger.closer)

	s.Len(logger.writers, 2)
}

func (s *LoggerTestSuite) TestNewLoggerCustomConfig() {
	logFile := filepath.Join(s.tempDir, "test-custom.log")

	logger, err := NewLogger(
		WithRotationEnabled(true),
		WithRotationFilename(logFile),
		WithRotationMaxSize(50),
		WithRotationMaxAge(7),
		WithRotationMaxBackups(3),
		WithRotationLocalTime(false),
		WithRotationCompress(false),
	)
	s.Require().NoError(err)
	s.NotNil(logger)

	rotator := logger.writers[1].(*lumberjack.Logger)
	s.Equal(logFile, rotator.Filename)
	s.Equal(50, rotator.MaxSize)
	s.Equal(3, rotator.MaxBackups)
	s.Equal(7, rotator.MaxAge)
	s.False(rotator.LocalTime)
	s.False(rotator.Compress)
}

func (s *LoggerTestSuite) TestLoggerOutputJSON() {
	var buf bytes.Buffer

	logFile := filepath.Join(s.tempDir, "test-json.log")
	logger, err := NewLogger(
		WithRotationFilename(logFile),
	)
	s.Require().NoError(err)
	s.NotNil(logger)

	logger.writers = []io.Writer{&buf}
	logger.Logger = slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))

	logger.Info("Test message", "key", "value")

	var logEntry map[string]any
	err = json.Unmarshal(buf.Bytes(), &logEntry)
	s.Require().NoError(err)

	s.Equal("INFO", logEntry["level"])
	s.Equal("Test message", logEntry["msg"])
	s.Equal("value", logEntry["key"])
	s.NotEmpty(logEntry["time"])
	s.NotEmpty(logEntry["source"])
}

func (s *LoggerTestSuite) TestLoggerOutputText() {
	var buf bytes.Buffer

	logFile := filepath.Join(s.tempDir, "test-text.log")
	logger, err := NewLogger(
		WithRotationFilename(logFile),
	)
	s.Require().NoError(err)
	s.NotNil(logger)

	logger.writers = []io.Writer{&buf}
	logger.Logger = slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))

	logger.Info("Test message", "key", "value")

	output := buf.String()
	s.Contains(output, "level=INFO")
	s.Contains(output, "msg=\"Test message\"")
	s.Contains(output, "key=value")
	s.Contains(output, "source=")
}

func (s *LoggerTestSuite) TestLoggerClose() {
	logFile := filepath.Join(s.tempDir, "test-close.log")
	logger, err := NewLogger(
		WithRotationFilename(logFile),
	)
	s.Require().NoError(err)

	err = logger.Close()
	s.Require().NoError(err)

	_, err = os.Stat(logFile)
	s.Require().NoError(err)
}

func (s *LoggerTestSuite) TestGetDefaultLogsPath() {
	s.Run("Dev mode", func() {
		version.SetVersion("dev")

		path := getDefaultLogsPath()

		s.Contains(path, filepath.Join("tmp", "logs", constant.AppName+".log"))
	})

	s.Run("Prod mode on Linux", func() {
		if runtime.GOOS == "linux" {
			version.SetVersion("1.2.3")

			// Test with XDG_STATE_HOME
			os.Setenv("XDG_STATE_HOME", "/test/xdg")
			path := getDefaultLogsPath()
			s.Equal(filepath.Join("/test/xdg", constant.AppName, "logs", constant.AppName+".log"), path)
			os.Unsetenv("XDG_STATE_HOME")

			// Test with home directory
			home, _ := os.UserHomeDir()
			path = getDefaultLogsPath()
			s.Contains(path, filepath.Join(home, ".local", "share", constant.AppName, "logs", constant.AppName+".log"))
		}
	})
}

func (s *LoggerTestSuite) TestLoggerDirectoryCreation() {
	logFile := filepath.Join(s.tempDir, "nested", "dir", "test.log")
	logger, err := NewLogger(
		WithRotationFilename(logFile),
	)
	s.Require().NoError(err)

	dir := filepath.Dir(logFile)
	_, err = os.Stat(dir)
	s.Require().NoError(err)

	err = logger.Close()
	s.Require().NoError(err)
}

func (s *LoggerTestSuite) TestLoggerErrorOnInvalidDirectory() {
	invalidDir := filepath.Join(s.tempDir, "invalid")
	err := os.WriteFile(invalidDir, []byte("test"), 0644)
	s.Require().NoError(err)

	logFile := filepath.Join(invalidDir, "test.log")
	logger, err := NewLogger(
		WithRotationFilename(logFile),
	)

	s.Require().Error(err)
	s.Nil(logger)
	s.Contains(err.Error(), "failed to create log directory")
}
