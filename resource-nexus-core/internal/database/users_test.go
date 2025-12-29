package database

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tbauriedel/resource-nexus-core/internal/config"
	database "github.com/tbauriedel/resource-nexus-core/internal/database/models"
	"github.com/tbauriedel/resource-nexus-core/internal/logging"
)

func TestGetUsers(t *testing.T) {
	d, mock, _ := sqlmock.New()
	defer d.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "password_hash", "is_admin"}).
		AddRow(1, "dummy", "foobar", false).
		AddRow(2, "dummy2", "foobar2", true)

	mock.ExpectQuery("SELECT id, name, password_hash, is_admin FROM users").WillReturnRows(rows)

	db := SqlDatabase{
		database: d,
		logger:   logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "warn"}),
	}

	ctx := context.Background()
	users, err := db.GetUsers(nil, ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(users) != 2 {
		t.Fatal("wrong number of users returned")
	}

	if users[0].Name != "dummy" || users[1].Name != "dummy2" {
		t.Fatal("wrong users returned")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

func TestGetUser(t *testing.T) {
	d, mock, _ := sqlmock.New()
	defer d.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "password_hash", "is_admin"}).
		AddRow(1, "dummy", "foobar", false)

	mock.ExpectQuery("SELECT id, name, password_hash, is_admin FROM users").WillReturnRows(rows)

	db := SqlDatabase{
		database: d,
		logger:   logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "warn"}),
	}

	ctx := context.Background()

	user, err := db.GetUser(Filter{Key: "name", Operator: "=", Value: "dummy"}, ctx)
	if err != nil {
		t.Fatal(err)
	}

	if user.Name != "dummy" {
		t.Fatal("wrong user returned")
	}
}

func TestInsertUser(t *testing.T) {
	d, mock, _ := sqlmock.New()
	defer d.Close()

	user := database.User{
		Name:         "dummy",
		PasswordHash: "foobar",
		IsAdmin:      false,
	}

	mock.ExpectExec(`INSERT INTO users \(name, password_hash, is_admin\) VALUES \(\$1, \$2, \$3\)`).
		WithArgs(user.Name, user.PasswordHash, user.IsAdmin).
		WillReturnResult(sqlmock.NewResult(1, 1))

	db := SqlDatabase{
		database: d,
		logger:   logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "warn"}),
	}

	ctx := context.TODO()
	res, err := db.InsertUser(ctx, user)
	if err != nil {
		t.Fatal(err)
	}

	if rows, _ := res.RowsAffected(); rows != 1 {
		t.Fatal("wrong number of rows affected")
	}
}

func TestGetUserPermissions(t *testing.T) {
	d, mock, _ := sqlmock.New()
	defer d.Close()

	query := fmt.Sprintf(`
		SELECT DISTINCT 
		    p.resource, p.action
		FROM %s u
		JOIN %s ug ON ug.user_id = u.id
		JOIN %s gp ON gp.group_id = ug.group_id
		JOIN %s p ON p.id = gp.permission_id
		`,
		TableNameUsers,
		TableNameUserGroups,
		TableNameGroupPermissions,
		TableNamePermissions,
	)

	rows := sqlmock.NewRows([]string{"resource", "action"}).
		AddRow("user", "get").
		AddRow("user", "create")

	mock.ExpectQuery(query).
		WithArgs("dummy").
		WillReturnRows(rows)

	db := SqlDatabase{
		database: d,
		logger:   logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "warn"}),
	}

	ctx := context.TODO()

	p, err := db.GetUserPermissions("dummy", ctx)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(p, []database.Permission{
		{Resource: "user", Action: "get"},
		{Resource: "user", Action: "create"}}) {
		t.Fatal("wrong permissions returned")
	}
}
