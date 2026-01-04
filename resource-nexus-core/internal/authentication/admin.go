package authentication

import (
	"context"
	"fmt"

	"github.com/tbauriedel/resource-nexus-core/internal/config"
	"github.com/tbauriedel/resource-nexus-core/internal/database"
)

// CreateAdminUser creates an admin user in the database.
//
// Provided password is hashed using the given hashing parameters.
func CreateAdminUser(password string, params config.HashingParams, db database.Database, ctx context.Context) error {
	filter := database.Filter{
		Key:      "name",
		Operator: "=",
		Value:    "admin",
	}

	_, err := db.GetUser(filter, ctx)
	if err == nil {
		return fmt.Errorf("admin user already exists")
	}

	user := database.User{
		Name:         "admin",
		IsAdmin:      true,
		PasswordHash: HashPasswordString(password, params),
	}

	// Insert user
	_, err = db.InsertUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to insert admin user: %w", err)
	}

	return nil
}
