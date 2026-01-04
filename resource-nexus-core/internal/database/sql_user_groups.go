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

	return getReferences(db, query, filter, ctx,
		func(rows *sql.Rows) (UserGroupReference, error) {
			var ref UserGroupReference

			err := rows.Scan(&ref.UserID, &ref.GroupID)
			if err != nil {
				return UserGroupReference{}, fmt.Errorf(
					"failed to scan group permission reference: %w",
					err,
				)
			}

			return ref, nil
		},
	)
}

// GetUserGroupReference returns a single user group reference that matches the given filter.
func (db *SqlDatabase) GetUserGroupReference(filter FilterExpr, ctx context.Context) (UserGroupReference, error) {
	userGroupRefs, err := db.GetUserGroupReferences(filter, ctx)
	if err != nil {
		return UserGroupReference{}, err
	}

	if !isSingleElement[UserGroupReference](userGroupRefs) {
		return UserGroupReference{},
			fmt.Errorf("not exactly 1 user group reference has been found with the filter %s", filter)
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
