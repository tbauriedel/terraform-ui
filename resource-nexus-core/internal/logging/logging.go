package logging

import (
	"log/slog"
	"os"

	"github.com/tbauriedel/resource-nexus-core/internal/config"
)

// Logger is an instance of slog.Logger.
type Logger struct {
	*slog.Logger
}

// mapLogLevel maps a string level to a slog.Level.
func mapLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// NewLoggerStdout creates a new logger that writes to stdout in JSON format.
func NewLoggerStdout(c config.Logger) *Logger {
	var l *slog.Logger

	opts := &slog.HandlerOptions{
		Level: mapLogLevel(c.Level),
	}

	l = slog.New(slog.NewJSONHandler(os.Stdout, opts))

	return &Logger{
		Logger: l,
	}
}

// NewLoggerFile creates a new logger that writes to a file in JSON format.
func NewLoggerFile(c config.Logger, f *os.File) *Logger {
	var l *slog.Logger

	// Set log options
	opts := &slog.HandlerOptions{
		Level: mapLogLevel(c.Level),
	}

	// Init new file logger in JSON format
	l = slog.New(slog.NewJSONHandler(f, opts))

	return &Logger{
		Logger: l,
	}
}
