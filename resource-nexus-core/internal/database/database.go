package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/tbauriedel/resource-nexus-core/internal/config"
	"github.com/tbauriedel/resource-nexus-core/internal/logging"
)

// Database interface for database operations.
// At the moment only PostgresSQL is supported.
// If more databases will be supported, they can implement the interface as well.
type Database interface { //nolint:interfacebloat
	TestConnection() error
	Close() error
	GetUsers(filter FilterExpr, ctx context.Context) ([]User, error)
	GetUser(filter FilterExpr, ctx context.Context) (User, error)
	GetUserPermissions(username string, ctx context.Context) ([]Permission, error)
	InsertUser(ctx context.Context, user User) (sql.Result, error)
	GetGroups(filter FilterExpr, ctx context.Context) ([]Group, error)
	GetGroup(filter FilterExpr, ctx context.Context) (Group, error)
	InsertGroup(ctx context.Context, group Group) (sql.Result, error)
	GetUserGroupReferences(filter FilterExpr, ctx context.Context) ([]UserGroupReference, error)
	GetUserGroupReference(filter FilterExpr, ctx context.Context) (UserGroupReference, error)
	InsertUserGroupReference(ctx context.Context, group UserGroupReference) (sql.Result, error)
	GetPermissions(filter FilterExpr, ctx context.Context) ([]Permission, error)
	GetPermission(filter FilterExpr, ctx context.Context) (Permission, error)
	GetGroupPermissions(filter FilterExpr, ctx context.Context) ([]GroupPermissionReference, error)
	GetGroupPermission(filter FilterExpr, ctx context.Context) (GroupPermissionReference, error)
	InsertGroupPermission(ctx context.Context, groupPermission GroupPermissionReference) (sql.Result, error)
}

type SqlDatabase struct {
	database *sql.DB
	logger   *logging.Logger
}

// sqlOpen is a mockable function to open a database connection.
var sqlOpen = sql.Open //nolint:gochecknoglobals

// NewDatabase creates a new SqlDatabase instance.
// The database connection is opened using the provided configuration.
// The database dsn will be built based on that config.
// PostgresSQL is used as the database engine.
// Returns an error if the connection fails.
func NewDatabase(conf config.Database, logger *logging.Logger) (*SqlDatabase, error) {
	db, err := sqlOpen("postgres", getDsn(conf))
	if err != nil {
		return &SqlDatabase{}, fmt.Errorf("failed to open database connection: %w", err)
	}

	return &SqlDatabase{
		database: db,
		logger:   logger,
	}, nil
}

// NewSqlDatabase returns a SqlDatabase instance with the given sql.DB connection.
//
// Dont use for prod instances. Use NewDatabase instead.
// This is only intended for testing purposes.
func NewSqlDatabase(db *sql.DB, l *logging.Logger) *SqlDatabase {
	return &SqlDatabase{
		database: db,
		logger:   l,
	}
}

// TestConnection tests the database connection by pinging the database.
// Returns an error if the ping fails.
func (db *SqlDatabase) TestConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.database.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// Close closes the database connection.
func (db *SqlDatabase) Close() error {
	db.logger.Info("closing database connection")

	return db.database.Close() //nolint:wrapcheck
}

// getDsn returns the DNS string for the database connection using postgresql.
//
// Format: postgres://username:password@localhost:5432/mydb?sslmode=verify-full
func getDsn(conf config.Database) string {
	// postgres://username:password@localhost:5432/mydb?sslmode=verify-full
	return fmt.Sprintf( //nolint:nosprintfhostport
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		conf.User,
		conf.Password,
		conf.Address,
		conf.Port,
		conf.Name,
		conf.TLSMode,
	)
}

type scanFn[T any] func(*sql.Rows) (T, error)

// getReferences executes a query and returns a list of items of type T.
func getReferences[T any](
	db *SqlDatabase,
	query string,
	filter FilterExpr,
	ctx context.Context,
	scan scanFn[T],
) ([]T, error) {
	rows, closeRows, err := db.Select(query, filter, ctx) //nolint:sqlclosecheck
	if err != nil {
		return nil, fmt.Errorf("failed to query references: %w", err)
	}
	defer closeRows()

	var result []T

	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reference: %w", err)
		}

		result = append(result, item)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %w", err)
	}

	return result, nil
}

// isSingleElement returns true if the given slice contains exactly one element.
func isSingleElement[T any](elements []T) bool {
	return len(elements) == 1
}
