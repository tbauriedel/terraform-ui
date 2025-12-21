package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tbauriedel/resource-nexus-core/internal/config"
	"github.com/tbauriedel/resource-nexus-core/internal/listener"
	"github.com/tbauriedel/resource-nexus-core/internal/logging"
	"github.com/tbauriedel/resource-nexus-core/internal/utils/fileutils"
	"github.com/tbauriedel/resource-nexus-core/internal/utils/netutils"
)

func main() {
	var (
		err        error
		conf       config.Config
		configPath string
		logger     *logging.Logger
	)

	// Set locale to C to avoid translations in command output
	_ = os.Setenv("LANG", "C")

	// Set slog defaults
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(handler))

	slog.Info("starting resource-nexus-core")

	// Add and parse flags
	flag.StringVar(&configPath, "config", "config.json", "Config file")

	flag.Parse()

	// Load server config
	if !fileutils.FileExists(configPath) {
		slog.Info("no config file found or provided. loading default configuration")

		conf = config.LoadDefaults()
	} else {
		slog.Info("loading config from file", "file", configPath)

		var err error

		conf, err = config.LoadFromJSONFile(configPath)
		if err != nil {
			slog.Error(err.Error())
			closeAndStop(nil, 1)
		}
	}

	slog.Info("configuration loaded")

	var f *os.File

	// init logger
	if conf.Logging.Type == "file" {
		slog.Debug("creating logger for type 'file'")

		// open logfile
		f, err = fileutils.OpenFile(conf.Logging.File)
		if err != nil {
			slog.Error(err.Error())
			closeAndStop(f, 1)
		}

		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		// create logger for type 'file'
		logger = logging.NewLoggerFile(conf.Logging, f)
	} else {
		slog.Debug("creating logger for type 'stdout'")

		// create logger for type 'stdout'
		logger = logging.NewLoggerStdout(conf.Logging)
	}

	// Get redacted config for initial log output
	redacted, err := json.Marshal(conf.GetConfigRedacted())
	if err != nil {
		logger.Error(fmt.Sprint("failed to marshal config", "error", err))

		closeAndStop(f, 1)
	}

	logger.Debug("starting with configuration", "config", redacted)

	//----- Listener -----//

	logger.Debug("initializing listener")

	// create new listener
	l := listener.NewListener(
		conf.Listener,
		logger,
		listener.WithMiddleWare(listener.MiddlewareLogging(logger)),
	)

	// Start listener in the background
	go func() {
		time.Sleep(2 * time.Second)

		logger.Info(fmt.Sprintf("starting listener on '%s'", conf.Listener.ListenAddr))

		err := l.Start()
		if err != nil {
			logger.Error(err.Error())
			closeAndStop(f, 1)
		}
	}()

	// Register interrupt signal handler
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	// Wait for the listener to start. shutdown after 10 seconds
	err = netutils.WaitForConnection(conf.Listener.ListenAddr, conf.Listener.TLSSkipVerify, 10*time.Second)
	if err != nil {
		logger.Error(fmt.Sprintf("waited 10 seconds for listener to start without success. shutting down. %s", err.Error()))
		closeAndStop(f, 1)
	}

	// send messages that resource-nexus-core is ready to server requests
	logger.Info(fmt.Sprintf("listener running on '%s'", conf.Listener.ListenAddr))
	logger.Info("resource-nexus-core ready. awaiting requests")

	// Wait for the interrupt signal to gracefully stop the listener
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("start shutting down resource-nexus-core")

	// shutdown context. Cancel will be called after 5 seconds
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	err = l.Shutdown(shutdownCtx)
	if err != nil {
		logger.Error(err.Error())
	}

	logger.Info("shut down resource-nexus-core")
}

// closeAndStop closes the logfile and exits the application with the given code.
func closeAndStop(logfile *os.File, code int) { //nolint:unparam
	if logfile != nil {
		_ = logfile.Close()
	}

	os.Exit(code)
}
