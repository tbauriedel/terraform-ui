package routes

import (
	"net/http"

	"github.com/tbauriedel/resource-nexus-core/internal/database"
)

type Routes struct {
	DB database.Database
}

type Route struct {
	Method      string
	Path        string
	HandlerFunc func(w http.ResponseWriter, r *http.Request)
}

func (routes *Routes) Get() []Route {
	return []Route{
		{
			Method:      http.MethodGet,
			Path:        "/system/health",
			HandlerFunc: routes.Health,
		},
	}
}
