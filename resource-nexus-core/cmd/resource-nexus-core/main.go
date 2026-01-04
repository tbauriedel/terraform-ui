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

	"github.com/tbauriedel/resource-nexus-core/internal/app"
	"github.com/tbauriedel/resource-nexus-core/internal/common/netutils"
	"github.com/tbauriedel/resource-nexus-core/internal/listener"
	"github.com/tbauriedel/resource-nexus-core/internal/logging"
)

func main() { //nolint:funlen,nolintlint,cyclop
	var (
		err        error
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

	// print shut down the message after all defer statements have been executed.
	defer func() {
		logger.Info("shut down resource-nexus-core")
	}()

	// Bootstrap application
	conf, err := app.LoadConfig(configPath)
	if err != nil {
		slog.Error(err.Error())
		app.Exit(nil, 1)
	}

	db, logger, logfile, err := app.Bootstrap(conf)
	if err != nil {
		slog.Error(err.Error())
		app.Exit(logfile, 1)
	}

	// close database connection on exit of main
	defer func() {
		err = db.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	//----- Listener -----//

	logger.Debug("initializing listener")

	// create new listener
	// several middlewares are added to the listener. (loaded from top to bottom)
	l := listener.NewListener(
		conf.Listener,
		logger,
		listener.WithMiddleWare(listener.MiddlewareRecovery(logger)), // fetch panics and recover
		listener.WithMiddleWare(listener.MiddlewareLogging(logger)),  // log all requests
		listener.WithMiddleWare(listener.MiddlewareGlobalRateLimiter( // rate limiting. global
			conf.Listener.GlobalRateLimitGeneration,
			conf.Listener.GlobalRateLimitBucketSize,
			logger,
		)),
		listener.WithMiddleWare(listener.MiddleWareIpRateLimiter( // rate limiting. ip-based
			conf.Listener.IpBasedRateLimitGeneration,
			conf.Listener.IpBasedRateLimitBucketSize,
			logger,
		)),
		listener.WithMiddleWare(listener.MiddlewareAuthentication(db, logger)), // validate user
	)

	// Add routes to the listener
	l.AddRoutesToListener(db, logger)

	// Start listener in the background
	go func() {
		err := l.Start()
		if err != nil {
			logger.Error(err.Error())
			app.Exit(logfile, 1)
		}
	}()

	// Register interrupt signal handler for the listener
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

	// Wait for the listener to start. shutdown after 10 seconds
	err = netutils.WaitForConnection(conf.Listener.ListenAddr, conf.Listener.TLSSkipVerify, 10*time.Second)
	if err != nil {
		logger.Error(fmt.Sprintf("waited 10 seconds for listener to start without success. shutting down. %s", err.Error()))
		app.Exit(logfile, 1)
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
