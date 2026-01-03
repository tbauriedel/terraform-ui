package database

import (
	"context"
	"database/sql"
)

// Select executes a query against the database.
//
// The final query is built out of the given query and filter expression. ctx is used to query the database.
//
// Returns selected rows and a function to close the rows.
func (db *SqlDatabase) Select(query string, filter FilterExpr, ctx context.Context) (*sql.Rows, func(), error) {
	// build where clause. returns empty string if no filter is given
	where, args, err := BuildWhere(filter)
	if err != nil {
		return nil, nil, err
	}

	// append where clause to query
	query += where

	db.logger.Debug("query database", "query", query, "args", args)

	// query database
	rows, err := db.database.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, func() {
			err = rows.Close()
			if err != nil {
				db.logger.Error("failed to close rows", "error", err)
			}
		}, err //nolint:wrapcheck
	}

	return rows, func() {
		err = rows.Close()
		if err != nil {
			db.logger.Error("failed to close rows", "error", err)
		}
	}, nil
}
