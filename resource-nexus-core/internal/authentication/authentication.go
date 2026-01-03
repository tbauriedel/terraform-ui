package authentication

import (
	"context"
	"fmt"
	"slices"

	"github.com/tbauriedel/resource-nexus-core/internal/database"
)

type User struct {
	ID              int
	Name            string
	PasswordHash    string
	Permissions     []string
	IsAdmin         bool
	IsAuthenticated bool
}

// LoadUser loads the user from the database.
// If the user does not exist, an error will be returned.
//
// Permissions will be loaded as well. To interact with the user, use the provided function of User.
func LoadUser(username string, db database.Database, ctx context.Context) (*User, error) {
	// prepare filter to query user by name
	filter := database.Filter{
		Key:      "name",
		Operator: "=",
		Value:    username,
	}

	// select user from the database
	res, err := db.GetUser(filter, ctx)
	if err != nil {
		return nil, fmt.Errorf("cant load user: %w", err)
	}

	user := &User{
		ID:              res.ID,
		Name:            res.Name,
		PasswordHash:    res.PasswordHash,
		IsAuthenticated: false,
		IsAdmin:         res.IsAdmin,
	}

	// Load Permissions from database
	perms, err := db.GetUserPermissions(user.Name, ctx)
	if err != nil {
		return nil, fmt.Errorf("cant load permissions for user %s: %w", user.Name, err)
	}

	for _, perm := range perms {
		user.Permissions = append(user.Permissions, BuildPermissionString(perm.Category, perm.Resource, perm.Action))
	}

	return user, nil
}

// Authenticate checks if the given password matches the stored hash.
//
// Extracts the parameters from the hash and compares the password with the encoded hash.
// If the password matches, true is returned, otherwise false.
func (user *User) Authenticate(password string) error {
	params, err := ExportParamsFromHash(user.PasswordHash)
	if err != nil {
		return fmt.Errorf("cant check if user is authenticated: %w", err)
	}

	encoded := HashPasswordString(password, params)

	if encoded == user.PasswordHash {
		user.IsAuthenticated = true

		return nil
	}

	return fmt.Errorf("password does not match")
}

// HasPermission checks if the user has the permission for the given resource and action.
//
// permission is in format "action:resource".
func (user *User) HasPermission(permission string) bool {
	if user.IsAdmin {
		return true
	}

	return slices.Contains(user.Permissions, permission)
}
