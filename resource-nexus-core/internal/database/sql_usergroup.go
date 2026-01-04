package database

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	TableNameUserGroups string = "user_groups"
)

// GetUserGroupReferences returns all user group references that match the given filter.
func (db *SqlDatabase) GetUserGroupReferences(filter FilterExpr, ctx context.Context) ([]UserGroupReference, error) {
	query := fmt.Sprintf("SELECT user_id, group_id FROM %s", TableNameUserGroups)

	rows, closeRows, err := db.Select(query, filter, ctx) //nolint:sqlclosecheck
	defer closeRows()

	if err != nil {
		return nil, err
	}

	var usersGroupRefs []UserGroupReference

	for rows.Next() {
		var userGroupRef UserGroupReference

		err = rows.Scan(&userGroupRef.UserID, &userGroupRef.GroupID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		usersGroupRefs = append(usersGroupRefs, userGroupRef)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %w", err)
	}

	return usersGroupRefs, nil
}

// GetUserGroupReference returns a single user group reference that matches the given filter.
func (db *SqlDatabase) GetUserGroupReference(filter FilterExpr, ctx context.Context) (UserGroupReference, error) {
	userGroupRefs, err := db.GetUserGroupReferences(filter, ctx)
	if err != nil {
		return UserGroupReference{}, err
	}

	if len(userGroupRefs) > 1 {
		return UserGroupReference{}, fmt.Errorf("found more than one user group reference with filter %s", filter)
	}

	if len(userGroupRefs) == 0 {
		return UserGroupReference{}, fmt.Errorf("no user group reference found with filter %s", filter)
	}

	return userGroupRefs[0], nil
}

// InsertUserGroupReference inserts a new user group reference into the database.
func (db *SqlDatabase) InsertUserGroupReference(ctx context.Context, group UserGroupReference) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, group_id) VALUES ($1, $2)", TableNameUserGroups)

	result, err := db.Insert(query, ctx, group.UserID, group.GroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user group reference: %w", err)
	}

	return result, nil
}
