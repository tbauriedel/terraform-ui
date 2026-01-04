package database

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	TableNameUsers string = "users"
)

// GetUsers returns all users from the database found by the given filter.
func (db *SqlDatabase) GetUsers(filter FilterExpr, ctx context.Context) ([]User, error) {
	query := fmt.Sprintf("SELECT id, name, password_hash, is_admin FROM %s", TableNameUsers)

	rows, closeRows, err := db.Select(query, filter, ctx) //nolint:sqlclosecheck
	if err != nil {
		return nil, err
	}

	defer closeRows()

	var users []User

	for rows.Next() {
		var user User

		err = rows.Scan(&user.ID, &user.Name, &user.PasswordHash, &user.IsAdmin)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %w", err)
	}

	return users, nil
}

// GetUser returns a single user from the database.
func (db *SqlDatabase) GetUser(filter FilterExpr, ctx context.Context) (User, error) {
	users, err := db.GetUsers(filter, ctx)
	if err != nil {
		return User{}, err
	}

	if len(users) > 1 {
		return User{}, fmt.Errorf("found more than one user with filter %s", filter)
	}

	if len(users) == 0 {
		return User{}, fmt.Errorf("no user found with filter %s", filter)
	}

	return users[0], nil
}

// InsertUser inserts a new user into the database.
func (db *SqlDatabase) InsertUser(ctx context.Context, user User) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, password_hash, is_admin) VALUES ($1, $2, $3)", TableNameUsers)

	result, err := db.Insert(query, ctx, user.Name, user.PasswordHash, user.IsAdmin)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return result, nil
}

// GetUserPermissions returns all permissions for a user.
// Permissions are joined from the user groups and the group permissions.
//
// The permissions are returned as a slice of database.Permission.
func (db *SqlDatabase) GetUserPermissions(username string, ctx context.Context) ([]Permission, error) {
	// SQL statement to query all permissions for a user
	query := fmt.Sprintf(`
		SELECT DISTINCT 
		    p.category, p.resource, p.action
		FROM %s u
		JOIN %s ug ON ug.user_id = u.id
		JOIN %s gp ON gp.group_id = ug.group_id
		JOIN %s p ON p.id = gp.permission_id
		`,
		TableNameUsers,
		TableNameUserGroups,
		TableNameGroupPermissions,
		TableNamePermissions,
	)

	// build filter for username
	filter := Filter{
		Key:      "name",
		Operator: "=",
		Value:    username,
	}

	// build where clause. is empty if no filter is given
	where, args, err := BuildWhere(filter)
	if err != nil {
		return nil, err
	}

	// append where clause to query
	query += where

	db.logger.Debug("query user permissions from database", "query", query, "args", args)

	rows, err := db.database.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query user permissions. query: %s: %w", query, err)
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			db.logger.Error("failed to close rows", "error", err)
		}
	}()

	var permissions []Permission

	for rows.Next() {
		var permission Permission

		err = rows.Scan(&permission.Category, &permission.Resource, &permission.Action)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row %w", err)
		}

		permissions = append(permissions, permission)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %w", err)
	}

	return permissions, nil
}
