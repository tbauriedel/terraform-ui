package database

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	TableNameUserGroups string = "user_groups"
)

// InsertUserGroupReference inserts a new user group reference into the database.
func (db *SqlDatabase) InsertUserGroupReference(ctx context.Context, group UserGroup) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (user_id, group_id) VALUES ($1, $2)", TableNameUserGroups)

	result, err := db.database.ExecContext(ctx, query, group.UserID, group.GroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user group reference: %w", err)
	}

	return result, nil
}
