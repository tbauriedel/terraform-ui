package authentication

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tbauriedel/resource-nexus-core/internal/config"
	"github.com/tbauriedel/resource-nexus-core/internal/database"
	"github.com/tbauriedel/resource-nexus-core/internal/logging"
)

func TestLoadUser(t *testing.T) {
	d, mock, _ := sqlmock.New()
	defer d.Close()

	// user rows
	rows := sqlmock.NewRows([]string{"id", "name", "password_hash", "is_admin"}).
		AddRow(1, "dummy", "foobar", false)

	mock.ExpectQuery("SELECT id, name, password_hash, is_admin FROM users").WillReturnRows(rows)

	// user permission rows
	query := fmt.Sprintf(`
		SELECT DISTINCT 
		    p.resource, p.action
		FROM %s u
		JOIN %s ug ON ug.user_id = u.id
		JOIN %s gp ON gp.group_id = ug.group_id
		JOIN %s p ON p.id = gp.permission_id
		`,
		"users",
		"user_groups",
		"group_permissions",
		"permissions",
	)

	rows = sqlmock.NewRows([]string{"resource", "action"}).
		AddRow("user", "get").
		AddRow("user", "create")

	mock.ExpectQuery(query).
		WithArgs("dummy").
		WillReturnRows(rows)

	l := logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "warn"})

	db := database.NewSqlDatabase(d, l)

	user, err := LoadUser("dummy", db, context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	expected := User{
		ID:           1,
		Name:         "dummy",
		PasswordHash: "foobar",
		Permissions:  []string{"get:user", "create:user"},
		IsAdmin:      false,
	}

	if !reflect.DeepEqual(user, &expected) {
		t.Fatalf("wrong user returned.\nactual: %v\nexpected: %v", user, &User{ID: 1, Name: "dummy", PasswordHash: "foobar", Permissions: []string{"get:user", "create:user"}, IsAdmin: false})
	}
}

func TestAuthenticate(t *testing.T) {
	// pass: foobar
	user := User{
		ID:           1,
		Name:         "foobar",
		PasswordHash: "$argon2id$v=19$m=65536,t=3,p=1$+kn21LcRetAkE7zObeS3xA$FzdfjLWlAiJbHLE+Rjm2hBMUmMb3TdmWQ7AMTtryYfk",
		Permissions:  []string{},
		IsAdmin:      false,
	}

	err := user.Authenticate("foobar")
	if err != nil {
		t.Fatal(err)
	}
}

func TestHasPermission(t *testing.T) {
	user := User{
		ID:          1,
		Name:        "dummy",
		Permissions: []string{"get:user", "create:user"},
		IsAdmin:     false,
	}

	if user.HasPermission("get:user") != true {
		t.Fatal("should have permission")
	}

	if user.HasPermission("create:cars") != false {
		t.Fatal("should not have permission")
	}

	admin := User{
		ID:          1,
		Name:        "admin",
		Permissions: []string{},
		IsAdmin:     true,
	}

	if admin.HasPermission("get:user") != true {
		t.Fatal("should have permission")
	}
}
