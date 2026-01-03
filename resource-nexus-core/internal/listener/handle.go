package listener

import (
	"fmt"
	"net/http"

	"github.com/tbauriedel/resource-nexus-core/internal/database"
	"github.com/tbauriedel/resource-nexus-core/internal/listener/routes"
)

// AddRoute adds a new route to the listener.
func (l *Listener) AddRoute(method, url string, handler http.HandlerFunc) {
	l.multiplexer.Handle(url, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			l.logger.Warn(fmt.Sprintf("method not allowed: %s %s", r.Method, r.URL.Path))

			return
		}

		MiddlewareAuthorization(l.logger)(handler).ServeHTTP(w, r)
	}))
}

// AddRoutesToListener adds all routes to the listener.
//
// Routes are defined in the 'routes' package.
func (l *Listener) AddRoutesToListener(db database.Database) {
	r := routes.Routes{
		DB: db,
	}

	for _, route := range r.Get() {
		l.AddRoute(route.Method, route.Path, route.HandlerFunc)
	}
}
