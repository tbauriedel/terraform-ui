package database

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	TableNameGroupPermissions string = "group_permissions"
)

// GetGroupPermissions returns all group permissions from the database based on the filter.
func (db *SqlDatabase) GetGroupPermissions(filter FilterExpr, ctx context.Context) ([]GroupPermissionReference, error) {
	query := fmt.Sprintf("SELECT group_id, permission_id FROM %s", TableNameGroupPermissions)

	return getReferences(db, query, filter, ctx,
		func(rows *sql.Rows) (GroupPermissionReference, error) {
			var ref GroupPermissionReference

			err := rows.Scan(&ref.GroupID, &ref.PermissionID)
			if err != nil {
				return GroupPermissionReference{}, fmt.Errorf(
					"failed to scan group permission reference: %w",
					err,
				)
			}

			return ref, nil
		},
	)
}

// GetGroupPermission returns a single group permission from the database based on the filter.
func (db *SqlDatabase) GetGroupPermission(filter FilterExpr, ctx context.Context) (GroupPermissionReference, error) {
	groupPermissions, err := db.GetGroupPermissions(filter, ctx)
	if err != nil {
		return GroupPermissionReference{}, err
	}

	if !isSingleElement[GroupPermissionReference](groupPermissions) {
		return GroupPermissionReference{},
			fmt.Errorf("not exactly 1 group permission reference has been found with the filter %s", filter)
	}

	return groupPermissions[0], nil
}

// InsertGroupPermission inserts a new group permission into the database.
func (db *SqlDatabase) InsertGroupPermission(
	ctx context.Context, groupPermission GroupPermissionReference,
) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (group_id, permission_id) VALUES ($1, $2)", TableNameGroupPermissions)

	result, err := db.Insert(query, ctx, groupPermission.GroupID, groupPermission.PermissionID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert group permission: %w", err)
	}

	return result, nil
}
