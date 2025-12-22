package listener

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/tbauriedel/resource-nexus-core/internal/logging"
	"golang.org/x/time/rate"
)

// Option for Listener.
type Option func(*Listener)

// Middleware is a function that wraps an http.Handler.
type Middleware func(http.Handler) http.Handler

// userRateLimiter holds the visitors rate limiters that are used by the MiddleWareIpRateLimiter.
type userRateLimiter struct {
	visitors            map[string]*rate.Limiter
	mu                  sync.Mutex
	rateLimitGeneration rate.Limit
	rateLimitBucketSize int
}

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
//	`r` is the rate at which new tokens are generated and stored to the bucket (tokens per second).
//	`b` is the maximum number of tokens the bucket can hold at once (burst capacity).
func MiddlewareGlobalRateLimiter(r rate.Limit, b int, logger *logging.Logger) Middleware {
	limiter := rate.NewLimiter(r, b)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				logger.Warn("global rate limit exceeded")

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// MiddleWareIpRateLimiter wraps the http.Handler with user rate limiting middleware.
//
// Rate limiting is token-based. Each request consumes a token.
// Tokens are generated and saved into a "bucket".
// When no token is available in the bucket, the request is rejected.
// Using this middleware, each user will have its own bucket based on the ip of the request.
//
//	`r` is the rate at which new tokens are generated and stored to the bucket (tokens per second).
//	`b` is the maximum number of tokens the bucket can hold at once (burst capacity).
func MiddleWareIpRateLimiter(r rate.Limit, b int, logger *logging.Logger) Middleware {
	l := &userRateLimiter{
		visitors:            make(map[string]*rate.Limiter),
		rateLimitGeneration: r,
		rateLimitBucketSize: b,
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// get ip from request
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				logger.Error(fmt.Sprintf("cant process user rate limiter: %s", err.Error()))

				return
			}

			// get limiter for the ip
			limiter := l.getIpRateLimit(ip)
			if !limiter.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				logger.Warn(fmt.Sprintf("ip-based rate limit exceeded: %s", r.RemoteAddr))

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getIpRateLimit returns the limiter for the given ip.
//
// If no limiter exists for the ip, a new one is created and stored in the visitors map.
func (r *userRateLimiter) getIpRateLimit(ip string) *rate.Limiter {
	// Lock mutex to avoid race conditions with more than one goroutine trying to access the map at the same time
	r.mu.Lock()
	defer r.mu.Unlock()

	// get limiter for the ip or create a new one if it does not exist
	limiter, exists := r.visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(r.rateLimitGeneration, r.rateLimitBucketSize)
		r.visitors[ip] = limiter
	}

	return limiter
}
