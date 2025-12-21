package listener

import (
	"fmt"
	"net/http"

	"github.com/tbauriedel/resource-nexus-core/internal/logging"
	"golang.org/x/time/rate"
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

// MiddlewareGlobalRateLimiter wraps the http.Handler with global rate limiting middleware.
//
// Rate limiting is token-based. Each request consumes a token.
// Tokens are generated and saved into a "bucket".
// When no token is available in the bucket, the request is rejected.
//
//	r is the rate at which new tokens are generated and stored to the bucket (tokens per second).
//	b is the maximum number of tokens the bucket can hold at once (burst capacity).
func MiddlewareGlobalRateLimiter(r rate.Limit, b int) Middleware {
	limiter := rate.NewLimiter(r, b)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
