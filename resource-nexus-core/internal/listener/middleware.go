package listener

import (
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/tbauriedel/resource-nexus-core/internal/authentication"
	"github.com/tbauriedel/resource-nexus-core/internal/database"
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

// MiddlewareRecovery wraps the http.Handler with recovery middleware.
//
// If a panic occurs during the request processing, server will log the panic and return a 500 Internal Server Error.
func MiddlewareRecovery(logger *logging.Logger) Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error(fmt.Sprintf("panic during request recovered: %v", err))
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
		})
	}
}

// MiddlewareLogging wraps the http.Handler with logging middleware.
func MiddlewareLogging(logger *logging.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, _, _ := r.BasicAuth()
			logger.Info(fmt.Sprintf(
				"new request from user '%s' [%s] %s %s %s",
				user, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent()))

			// hand over to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// MiddlewareAuthentication wraps the http.Handler with authentication middleware.
//
// The middleware expects a valid username and password in the BasicAuth header of the request.
// The provided password is validated against the stored hash in the database.
// If the authentication fails, the request is rejected with a 401 Unauthorized status code.
// The http response, no proper error message is returned to prevent leaking information about the existence of users.
// Details about the failed authentication are logged inside the server log.
func MiddlewareAuthentication(db database.Database, logger *logging.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, _ := r.BasicAuth()

			// if no basic auth header is provided, reject the request immediately
			if username == "" || password == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				logger.Warn("authentication failed: no username or password provided")

				return
			}

			// load user from database by given name.
			// unauthorized if no user can be found
			storedUser, err := authentication.LoadUser(username, db, r.Context())
			if err != nil {
				logger.Warn(fmt.Sprintf("authentication failed: %s", err.Error()))
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

				return
			}

			// validate the provided password against the stored hash
			err = storedUser.Authenticate(password)
			if err != nil {
				logger.Warn(fmt.Sprintf("authentication for user '%s' failed: %s", username, err.Error()))
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)

				return
			}

			logger.Info(fmt.Sprintf("authentication for user '%s' successful", username))

			// hand over to the next handler
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

			// hand over to the next handler
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

			// hand over to the next handler
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
