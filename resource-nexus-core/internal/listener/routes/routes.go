package routes

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tbauriedel/resource-nexus-core/internal/database"
	"github.com/tbauriedel/resource-nexus-core/internal/logging"
)

type Routes struct {
	DB     database.Database
	Logger *logging.Logger
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
		{
			Method:      http.MethodPost,
			Path:        "/auth/user/add",
			HandlerFunc: routes.UserAdd,
		},
		{
			Method:      http.MethodPost,
			Path:        "/auth/group/add",
			HandlerFunc: routes.GroupAdd,
		},
		{
			Method:      http.MethodPost,
			Path:        "/auth/usergroup/add",
			HandlerFunc: routes.AddUserToGroup,
		},
		{
			Method:      http.MethodPost,
			Path:        "/auth/grouppermission/add",
			HandlerFunc: routes.AddPermissionToGroup,
		},
	}
}

// BuildResponseMessage builds a response message for the client.
func BuildResponseMessage(message string) string {
	return fmt.Sprintf("{\"message\":\"%s\"}", message)
}

// decodeJson decodes a json request body into a struct.
//
// Make sure to provide a valid type when calling this function!
func decodeJson[T any](r *http.Request) (T, error) { //nolint:ireturn
	defer r.Body.Close()

	var obj T

	err := json.NewDecoder(r.Body).Decode(&obj)

	return obj, err //nolint:wrapcheck
}

// addEntity adds a new entity to the database.
//
// It takes the current request and response writer.
// The entity that will be added needs to be provided as the entity parameter.
// The filter is used to check if the entity already exists inside the getFunc function.
// The getFunc and insertFunc are used to query and insert the entity into the database.
func addEntity(
	w http.ResponseWriter,
	r *http.Request,
	entity any,
	filter database.Filter,
	getFunc func(database.FilterExpr, context.Context) (any, error),
	insertFunc func(context.Context, any) (sql.Result, error),
) error {
	// validate if entity with the same name already exists
	// takes the given function to query the database with the provided filter
	_, err := getFunc(filter, r.Context())
	if err == nil {
		http.Error(w,
			BuildResponseMessage("entity with the same name already exists"),
			http.StatusBadRequest)

		return fmt.Errorf("entity with the same name already exists")
	}

	// inserts the entity into the database with the given insert function
	result, err := insertFunc(r.Context(), entity)
	if err != nil {
		http.Error(w,
			BuildResponseMessage("failed to insert user. check logs for more details"),
			http.StatusInternalServerError,
		)

		return fmt.Errorf("failed to insert entity: %w", err)
	}

	// checks if the insert operation was successful
	rows, _ := result.RowsAffected()
	if rows != 1 {
		http.Error(w,
			BuildResponseMessage("error while inserting user. check logs for more details"),
			http.StatusInternalServerError,
		)

		return fmt.Errorf("error while inserting entity: %w", err)
	}

	// sends response to client
	_, _ = w.Write([]byte(BuildResponseMessage("entity created successfully")))

	return nil
}
