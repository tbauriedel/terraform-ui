package routes

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/tbauriedel/resource-nexus-core/internal/database"
)

func (routes *Routes) UserAdd(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	user, err := decodeJson[database.User](r)
	if err != nil {
		http.Error(w,
			BuildResponseMessage("invalid json"),
			http.StatusBadRequest)
		routes.Logger.Error("failed to decode user from body", "error", err)

		return
	}

	err = addEntity(
		w, r,
		user,
		database.Filter{Key: "name", Operator: "=", Value: user.Name},
		func(filter database.FilterExpr, ctx context.Context) (any, error) {
			return routes.DB.GetUser(filter, ctx)
		},
		func(ctx context.Context, entity any) (sql.Result, error) {
			return routes.DB.InsertUser(ctx, user)
		},
	)

	routes.Logger.Error("failed to add user", "error", err)
}

// GroupAdd adds a new group to the database.
func (routes *Routes) GroupAdd(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	group, err := decodeJson[database.Group](r)
	if err != nil {
		http.Error(w,
			BuildResponseMessage("invalid json"),
			http.StatusBadRequest,
		)
		routes.Logger.Error("failed to decode group from body", "error", err)

		return
	}

	err = addEntity(
		w, r,
		group,
		database.Filter{Key: "name", Operator: "=", Value: group.Name},
		func(filter database.FilterExpr, ctx context.Context) (any, error) {
			return routes.DB.GetGroup(filter, ctx)
		},
		func(ctx context.Context, entity any) (sql.Result, error) {
			return routes.DB.InsertGroup(ctx, group)
		},
	)
	if err != nil {
		routes.Logger.Error("failed to add group", "error", err)
	}
}
