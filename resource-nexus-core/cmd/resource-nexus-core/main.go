package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tbauriedel/resource-nexus-core/internal/common/fileutils"
	"github.com/tbauriedel/resource-nexus-core/internal/common/netutils"
	"github.com/tbauriedel/resource-nexus-core/internal/config"
	"github.com/tbauriedel/resource-nexus-core/internal/database"
	"github.com/tbauriedel/resource-nexus-core/internal/listener"
	"github.com/tbauriedel/resource-nexus-core/internal/logging"
)

func main() { //nolint:funlen,nolintlint,cyclop
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
			exit(nil, 1)
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
			exit(f, 1)
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

	// Print redacted config
	logger.Debug("starting with configuration", "config", conf.GetConfigRedacted())

	// print shut down the message after all defer statements have been executed.
	defer func() {
		logger.Info("shut down resource-nexus-core")
	}()

	//----- Database -----//

	logger.Info("initializing database connection")

	var db database.Database

	// create the database connection
	db, err = database.NewDatabase(conf.Database, logger)
	if err != nil {
		logger.Error(err.Error())
		exit(f, 1)
	}

	// close database connection on exit of main
	defer func() {
		err = db.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	// test database connection
	err = db.TestConnection()
	if err != nil {
		logger.Error(err.Error())
		exit(f, 1)
	}

	logger.Info("database connection established and tested successfully")

	//----- Listener -----//

	logger.Debug("initializing listener")

	// create new listener
	// several middlewares are added to the listener. (loaded from top to bottom)
	// - MiddleWareRecovery: recovers from panics and logs the error
	// - MiddleWareLogging: logs the request and response
	// - MiddleWareGlobalRateLimiter: limits the number of requests per second globally
	// - MiddleWareIpRateLimiter: limits the number of requests per second per ip
	// - MiddlewareAuthentication: authenticates requests
	l := listener.NewListener(
		conf.Listener,
		logger,
		listener.WithMiddleWare(listener.MiddlewareRecovery(logger)),
		listener.WithMiddleWare(listener.MiddlewareLogging(logger)),
		listener.WithMiddleWare(listener.MiddlewareGlobalRateLimiter(
			conf.Listener.GlobalRateLimitGeneration,
			conf.Listener.GlobalRateLimitBucketSize,
			logger,
		)),
		listener.WithMiddleWare(listener.MiddleWareIpRateLimiter(
			conf.Listener.IpBasedRateLimitGeneration,
			conf.Listener.IpBasedRateLimitBucketSize,
			logger,
		)),
		listener.WithMiddleWare(listener.MiddlewareAuthentication(db, logger)),
	)

	// Start listener in the background
	go func() {
		err := l.Start()
		if err != nil {
			logger.Error(err.Error())
			exit(f, 1)
		}
	}()

	// Register interrupt signal handler for the listener
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	// Wait for the listener to start. shutdown after 10 seconds
	err = netutils.WaitForConnection(conf.Listener.ListenAddr, conf.Listener.TLSSkipVerify, 10*time.Second)
	if err != nil {
		logger.Error(fmt.Sprintf("waited 10 seconds for listener to start without success. shutting down. %s", err.Error()))
		exit(f, 1)
	}

	// send messages that resource-nexus-core is ready to server requests
	logger.Info(fmt.Sprintf("listener running on '%s'", conf.Listener.ListenAddr))
	logger.Info("resource-nexus-core ready. awaiting requests")

	// ----- Shutdown ----- //

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

	logger.Debug("listener stopped")
}

// exit closes the logfile and exits the application with the given code.
func exit(logfile *os.File, code int) { //nolint:unparam
	if logfile != nil {
		_ = logfile.Close()
	}

	os.Exit(code)
}
