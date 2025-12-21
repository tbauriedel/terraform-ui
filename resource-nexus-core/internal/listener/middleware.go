package listener

import (
	"fmt"
	"net/http"

	"github.com/tbauriedel/resource-nexus-core/internal/logging"
)

// Option for Listener.
type Option func(*Listener)

// Middleware is a function that wraps an http.Handler.
type Middleware func(http.Handler) http.Handler

// WithMiddleWare adds middleware to the listener.
func WithMiddleWare(m Middleware) Option {
	return func(l *Listener) {
		l.middlewares = append(l.middlewares, m)
	}
}

// MiddlewareLogging wraps the http.Handler with logging middleware.
func MiddlewareLogging(logger *logging.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(fmt.Sprintf("new request: [%s] %s %s %s", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent()))
			next.ServeHTTP(w, r)
		})
	}
}
