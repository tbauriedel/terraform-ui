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

	"github.com/tbauriedel/terraform-ui-core/internal/config"
	"github.com/tbauriedel/terraform-ui-core/internal/listener"
	"github.com/tbauriedel/terraform-ui-core/internal/logging"
	"github.com/tbauriedel/terraform-ui-core/internal/utils/fileutils"
	"github.com/tbauriedel/terraform-ui-core/internal/utils/netutils"
)

var (
	configPath string
	conf       config.Config
	logger     *logging.Logger
	err        error
)

func init() {
	// Set locale to C to avoid translations in command output
	_ = os.Setenv("LANG", "C")

	// Set slog defaults
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	slog.SetDefault(slog.New(handler))

	slog.Info("starting terraform-ui-core")

	// Add and parse flags
	flag.StringVar(&configPath, "config", "config.json", "Config file")

	flag.Parse()

	// Load server config
	if !fileutils.FileExists(configPath) {
		slog.Info("no config file found or provided. loading default configuration")

		conf = config.LoadDefaults()
	} else {
		slog.Info("loading config from file", "file", configPath)

		conf, err = config.LoadFromJSONFile(configPath)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}

	slog.Info("configuration loaded")
}

func main() {
	// init logger
	if conf.Logging.Type == "file" {
		slog.Debug("creating logger for type 'file'")

		// open logfile
		f, err := fileutils.OpenFile(conf.Logging.File)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}

		// close the logfile when main is done
		defer f.Close()

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
		os.Exit(1)
	}

	logger.Debug("starting with configuration", "config", redacted)

	//----- Listener -----//

	logger.Debug("initializing listener")

	// create new listener
	l := listener.NewListener(conf.Listener, context.Background(), logger)

	// Start listener in the background
	go func() {
		time.Sleep(2 * time.Second)

		logger.Info(fmt.Sprintf("starting listener on '%s'", conf.Listener.ListenAddr))
		err := l.Start()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}()

	// Register interrupt signal handler
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	// Wait for the listener to start. shutdown after 10 seconds
	err = netutils.WaitForConnection(conf.Listener.ListenAddr, 10*time.Second)
	if err != nil {
		logger.Error(fmt.Sprintf("waited 10 seconds for listener to start without success. shutting down. %s", err.Error()))
		os.Exit(1)
	}

	// send messages that terraform-ui-core is ready to server requests
	logger.Info(fmt.Sprintf("listener running on '%s'", conf.Listener.ListenAddr))
	logger.Info("terraform-ui-core ready. awaiting requests")

	// Wait for the interrupt signal to gracefully stop the listener
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("start shutting down terraform-ui-core")

	// shutdown context. Cancel will be called after 5 seconds
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	err = l.Shutdown(shutdownCtx)
	if err != nil {
		logger.Error(err.Error())
	}

	logger.Info("shut down terraform-ui-core")

	return
}
