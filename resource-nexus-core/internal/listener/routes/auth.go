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

// AddUserToGroup adds a user to a group based on the group name and username.
func (routes *Routes) AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userGroupRef, err := decodeJson[database.UserGroupReference](r)
	if err != nil {
		http.Error(w,
			BuildResponseMessage("invalid json"),
			http.StatusBadRequest,
		)
		routes.Logger.Error("failed to decode user group reference from body", "error", err)
	}

	// get user by name
	user, err := routes.DB.GetUser(database.Filter{
		Key:      "name",
		Operator: "=",
		Value:    userGroupRef.Username,
	}, r.Context())
	if err != nil {
		http.Error(w,
			BuildResponseMessage("user not found"),
			http.StatusBadRequest,
		)
		routes.Logger.Error("failed to get user", "error", err)

		return
	}

	userGroupRef.UserID = user.ID

	// get group by name
	group, err := routes.DB.GetGroup(database.Filter{
		Key:      "name",
		Operator: "=",
		Value:    userGroupRef.GroupName,
	}, r.Context())
	if err != nil {
		http.Error(w,
			BuildResponseMessage("group not found"),
			http.StatusBadRequest,
		)
		routes.Logger.Error("failed to get group", "error", err)

		return
	}

	userGroupRef.GroupID = group.ID

	// check if user is already in group
	_, err = routes.DB.GetUserGroupReference(database.LogicalFilter{
		Operator: "AND",
		Filters: []database.FilterExpr{
			database.Filter{
				Key:      "group_id",
				Operator: "=",
				Value:    group.ID,
			},
			database.Filter{
				Key:      "user_id",
				Operator: "=",
				Value:    user.ID,
			},
		},
	}, r.Context())
	if err == nil {
		http.Error(w,
			BuildResponseMessage("user already in group"),
			http.StatusBadRequest,
		)
		routes.Logger.Error("user already in group", "error", err)

		return
	}
	
	result, err := routes.DB.InsertUserGroupReference(r.Context(), userGroupRef)
	if err != nil {
		http.Error(w,
			BuildResponseMessage("failed to insert user group reference"),
			http.StatusInternalServerError,
		)
		routes.Logger.Error("failed to insert user group reference", "error", err)
	}

	rows, _ := result.RowsAffected()
	if rows != 1 {
		http.Error(w,
			BuildResponseMessage("error while inserting user. check logs for more details"),
			http.StatusInternalServerError,
		)
		routes.Logger.Error("error while inserting user group reference", "error", err)

		return
	}

	_, _ = w.Write([]byte(BuildResponseMessage("user group reference added")))
}
