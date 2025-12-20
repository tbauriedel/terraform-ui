package listener

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"github.com/tbauriedel/terraform-ui-core/internal/config"
	"github.com/tbauriedel/terraform-ui-core/internal/logging"
)

type Listener struct {
	logger      *logging.Logger
	config      config.Listener
	context     context.Context
	multiplexer http.Handler
	server      *http.Server
}

// NewListener retunrs a new and empty Listener
func NewListener(cfg config.Listener, ctx context.Context, logger *logging.Logger) Listener {
	// Create a new and empty multiplexer
	multiplexer := http.NewServeMux()

	return Listener{
		logger:      logger,
		config:      cfg,
		context:     ctx,
		multiplexer: multiplexer,
	}
}

// Start starts the listener in the foreground.
// Config will be loaded from the Listener config
//
// If TLS is enabled, the server will be started with TLS encryption. If not, the server will be started without TLS encryption
func (l *Listener) Start() error {
	server := &http.Server{
		Addr:        l.config.ListenAddr,
		Handler:     l.multiplexer,
		ReadTimeout: l.config.ReadTimeout,
		IdleTimeout: l.config.IdleTimeout,
	}

	// add the server to the listener to have it available in the Shutdown function
	l.server = server

	// if TLS is not enabled, run the server without TLS
	if !l.config.TlsEnabled {
		l.logger.Info("tls is disabled. running listener without TLS encryption")

		// run the server without TLS encryption
		err := l.run(server)
		if err != nil {
			return err
		}

		return nil
	}

	// Prepare TLS server
	// Load and parse and key file
	l.logger.Debug(fmt.Sprintf("loading tls certificate and key from files. crt: %s, key: %s", l.config.TlsCertPath, l.config.TlsKeyPath))

	cert, err := tls.LoadX509KeyPair(l.config.TlsCertPath, l.config.TlsKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load TLS certificate: %w", err)
	}

	// Load TLS certificate and key
	server.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: l.config.TlsSkipVerify}

	// Run the server with TLS encryption
	l.logger.Info("tls is enabled. running listener with tls encryption")

	// start secure listener server
	err = l.runSecure(server)
	if err != nil {
		return err
	}

	return nil
}

// run starts the server without TLS encryption in the foreground
func (l *Listener) run(s *http.Server) error {
	if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) && err != nil {
		return fmt.Errorf("failed to start insecure listener: %w", err)
	}

	return nil
}

// runSecure starts the server with TLS encryption in the foreground
func (l *Listener) runSecure(s *http.Server) error {
	// Cert and key file already loaded in the server. No need to load them again
	if err := s.ListenAndServeTLS("", ""); !errors.Is(err, http.ErrServerClosed) && err != nil {
		return fmt.Errorf("failed to start secure listener: %w", err)
	}

	return nil
}

// Shutdown shuts down the listener gracefully
func (l *Listener) Shutdown(ctx context.Context) error {
	l.logger.Debug("shutdown for listener requested")

	err := l.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("failed to shutdown listener: %w", err)
	}

	return nil
}
