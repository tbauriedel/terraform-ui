package routes

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

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

	_, err = routes.DB.InsertUserGroupReference(r.Context(), userGroupRef)
	if err != nil {
		http.Error(w,
			BuildResponseMessage("failed to insert user group reference"),
			http.StatusInternalServerError,
		)
		routes.Logger.Error("failed to insert user group reference", "error", err)
	}

	_, _ = w.Write([]byte(BuildResponseMessage("user group reference added")))
}

func (routes *Routes) AddPermissionToGroup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	
	// decode permission group from json body
	permissionGroup, err := decodeJson[database.GroupPermissionReference](r)
	if err != nil {
		http.Error(w,
			BuildResponseMessage("invalid json"),
			http.StatusBadRequest,
		)
		routes.Logger.Error("failed to decode permission group from body", "error", err)

		return
	}

	// check if group exists. get group by name
	group, err := routes.DB.GetGroup(database.Filter{
		Key:      "name",
		Operator: "=",
		Value:    permissionGroup.GroupName,
	}, r.Context())
	if err != nil {
		http.Error(w,
			BuildResponseMessage("group not found"),
			http.StatusBadRequest,
		)
		routes.Logger.Error("failed to get group", "error", err)

		return
	}

	permissionGroup.GroupID = group.ID

	// check if permission exists
	// split permission string into category, resource and action
	splitted := strings.Split(permissionGroup.Permission, ":")

	// prepare filter
	filter := database.LogicalFilter{
		Operator: "AND",
		Filters: []database.FilterExpr{
			database.Filter{
				Key:      "category",
				Operator: "=",
				Value:    splitted[0],
			},
			database.Filter{
				Key:      "resource",
				Operator: "=",
				Value:    splitted[1],
			},
			database.Filter{
				Key:      "action",
				Operator: "=",
				Value:    splitted[2],
			},
		},
	}

	permission, err := routes.DB.GetPermission(filter, r.Context())
	if err != nil {
		http.Error(w,
			BuildResponseMessage("permission not found"),
			http.StatusBadRequest,
		)
		routes.Logger.Error("failed to get permission", "error", err)

		return
	}

	permissionGroup.PermissionID = permission.ID

	// check if group already has permission
	_, err = routes.DB.GetGroupPermission(database.LogicalFilter{
		Operator: "AND",
		Filters: []database.FilterExpr{
			database.Filter{
				Key:      "group_id",
				Operator: "=",
				Value:    group.ID,
			},
			database.Filter{
				Key:      "permission_id",
				Operator: "=",
				Value:    permission.ID,
			},
		},
	}, r.Context())
	if err == nil {
		http.Error(w,
			BuildResponseMessage("permission already assigned to group"),
			http.StatusBadRequest,
		)
		routes.Logger.Error("permission already assigned to group", "error", err)

		return
	}

	// add permission to group
	_, err = routes.DB.InsertGroupPermission(r.Context(), permissionGroup)
	if err != nil {
		http.Error(w,
			BuildResponseMessage("failed to insert permission group reference"),
			http.StatusInternalServerError,
		)
		routes.Logger.Error("failed to insert permission group reference", "error", err)

		return
	}

	_, _ = w.Write([]byte(BuildResponseMessage("permission group reference added")))
}
