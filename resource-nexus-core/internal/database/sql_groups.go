package database

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	TableNameGroups string = "groups"
)

// GetGroups returns all groups from the database based on the filter.
func (db *SqlDatabase) GetGroups(filter FilterExpr, ctx context.Context) ([]Group, error) {
	query := fmt.Sprintf("SELECT id, name FROM %s", TableNameGroups)

	// query database
	rows, closeRows, err := db.Select(query, filter, ctx) //nolint:sqlclosecheck
	defer closeRows()

	if err != nil {
		return nil, err
	}

	var groups []Group

	for rows.Next() {
		var group Group

		err = rows.Scan(&group.ID, &group.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		groups = append(groups, group)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %w", err)
	}

	return groups, nil
}

// GetGroup returns a single group from the database based on the filter.
func (db *SqlDatabase) GetGroup(filter FilterExpr, ctx context.Context) (Group, error) {
	groups, err := db.GetGroups(filter, ctx)
	if err != nil {
		return Group{}, err
	}

	if len(groups) > 1 {
		return Group{}, fmt.Errorf("found more than one group with filter %s", filter)
	}

	if len(groups) == 0 {
		return Group{}, fmt.Errorf("no group found with filter %s", filter)
	}

	return groups[0], nil
}

// InsertGroup inserts a new group into the database.
func (db *SqlDatabase) InsertGroup(ctx context.Context, group Group) (sql.Result, error) {
	query := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1)", TableNameGroups)

	result, err := db.Insert(query, ctx, group.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to insert group: %w", err)
	}

	return result, nil
}
