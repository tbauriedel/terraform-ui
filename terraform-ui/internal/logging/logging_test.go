package logging

import (
	"log/slog"
	"os"
	"testing"

	"github.com/tbauriedel/terraform-ui/internal/config"
)

func TestMapLogLevel(t *testing.T) {
	types := map[string]slog.Level{
		"debug":   slog.LevelDebug,
		"warn":    slog.LevelWarn,
		"error":   slog.LevelError,
		"default": slog.LevelInfo, // default is info
	}

	for x, y := range types {
		level := mapLogLevel(x)
		if level != y {
			t.Fatalf("Provided: %s. result expected: %s, result actual: %s", x, y, level)
		}
	}
}

func TestNewLoggerStdout(t *testing.T) {
	c := config.LoadDefaults()
	c.Logging.Type = "stdout"

	l := NewLoggerStdout(c.Logging)

	l.Info("Test")
}

func TestNewLoggerFile(t *testing.T) {
	c := config.Config{
		Logging: config.Logger{
			Type: "file",
			File: "../../test/testdata/test.log",
		},
	}

	f, err := os.OpenFile(c.Logging.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatal(err)
	}

	l := NewLoggerFile(c.Logging, f)

	l.Info("Test")
}
