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
type Database interface {
	TestConnection() error
	Close() error
	GetUsers(filter FilterExpr, ctx context.Context) ([]User, error)
	GetUser(filter FilterExpr, ctx context.Context) (User, error)
	GetUserPermissions(username string, ctx context.Context) ([]Permission, error)
	InsertUser(ctx context.Context, user User) (sql.Result, error)
	GetGroups(filter FilterExpr, ctx context.Context) ([]Group, error)
	GetGroup(filter FilterExpr, ctx context.Context) (Group, error)
	InsertGroup(ctx context.Context, group Group) (sql.Result, error)
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
