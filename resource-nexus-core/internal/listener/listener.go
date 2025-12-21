package listener

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"github.com/tbauriedel/resource-nexus-core/internal/config"
	"github.com/tbauriedel/resource-nexus-core/internal/logging"
)

type Listener struct {
	logger      *logging.Logger
	config      config.Listener
	multiplexer http.Handler
	server      *http.Server
	middlewares []Middleware
}

// NewListener returns a new and empty Listener
// The listener will be initialized with the given options
//
// Added middlewares will be applied to the multiplexer.
// The first one added, will be called first, and the last one added, will be called last.
func NewListener(cfg config.Listener, logger *logging.Logger, opts ...Option) *Listener {
	l := &Listener{
		logger: logger,
		config: cfg,
	}

	// apply options to listener
	for _, opt := range opts {
		opt(l)
	}

	if l.multiplexer == nil {
		l.multiplexer = http.NewServeMux()
	}

	// apply middlewares to mux
	l.applyMiddlewares()

	return l
}

// Start starts the listener in the foreground.
// Config will be loaded from the Listener config
//
// If TLS is enabled, the server will be started with TLS encryption.
// If not, the server will be started without TLS encryption.
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
	if !l.config.TLSEnabled {
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
	l.logger.Debug(fmt.Sprintf("loading tls certificate and key from files. crt: %s, key: %s",
		l.config.TLSCertPath,
		l.config.TLSKeyPath,
	))

	cert, err := tls.LoadX509KeyPair(l.config.TLSCertPath, l.config.TLSKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load TLS certificate: %w", err)
	}

	// Load TLS certificate and key
	server.TLSConfig = &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: l.config.TLSSkipVerify, //nolint:gosec
	}

	// Run the server with TLS encryption
	l.logger.Info("tls is enabled. running listener with tls encryption")

	// start secure listener server
	err = l.runSecure(server)
	if err != nil {
		return err
	}

	return nil
}

// Shutdown shuts down the listener gracefully.
func (l *Listener) Shutdown(ctx context.Context) error {
	l.logger.Debug("shutdown for listener requested")

	err := l.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("failed to shutdown listener: %w", err)
	}

	return nil
}

// applyMiddlewares applies the given middlewares to the multiplexer.
// Middlewares are applied in reverse order, so the first middleware will be applied last
//
// Listener.multiplexer needs to be set before calling this function.
func (l *Listener) applyMiddlewares() {
	for i := len(l.middlewares) - 1; i >= 0; i-- {
		l.multiplexer = l.middlewares[i](l.multiplexer)
	}
}

// run starts the server without TLS encryption in the foreground.
func (l *Listener) run(s *http.Server) error {
	err := s.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) && err != nil {
		return fmt.Errorf("failed to start insecure listener: %w", err)
	}

	return nil
}

// runSecure starts the server with TLS encryption in the foreground.
func (l *Listener) runSecure(s *http.Server) error {
	// Cert and key file already loaded in the server. No need to load them again
	err := s.ListenAndServeTLS("", "")
	if !errors.Is(err, http.ErrServerClosed) && err != nil {
		return fmt.Errorf("failed to start secure listener: %w", err)
	}

	return nil
}
