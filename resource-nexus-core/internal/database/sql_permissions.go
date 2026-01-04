package database

import (
	"context"
	"fmt"
)

const (
	TableNamePermissions string = "permissions"
)

// GetPermissions returns all permissions from the database based on the filter.
func (db *SqlDatabase) GetPermissions(filter FilterExpr, config context.Context) ([]Permission, error) {
	query := fmt.Sprintf("SELECT id, category, resource, action FROM %s", TableNamePermissions)

	rows, closeRows, err := db.Select(query, filter, config) //nolint:sqlclosecheck
	if err != nil {
		return nil, err
	}

	defer closeRows()

	var permissions []Permission

	for rows.Next() {
		var permission Permission

		err = rows.Scan(&permission.ID, &permission.Category, &permission.Resource, &permission.Action)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}

		permissions = append(permissions, permission)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %w", err)
	}

	return permissions, nil
}

// GetPermission returns a single permission from the database based on the filter.
func (db *SqlDatabase) GetPermission(filter FilterExpr, ctx context.Context) (Permission, error) {
	permissions, err := db.GetPermissions(filter, ctx)
	if err != nil {
		return Permission{}, err
	}

	if len(permissions) != 1 {
		return Permission{}, fmt.Errorf("found more than one permission with filter %s", filter)
	}

	if len(permissions) == 0 {
		return Permission{}, fmt.Errorf("no permission found with filter %s", filter)
	}

	return permissions[0], nil
}
