package database

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tbauriedel/resource-nexus-core/internal/config"
	"github.com/tbauriedel/resource-nexus-core/internal/logging"
)

func TestGetPermissions(t *testing.T) {
	d, mock, _ := sqlmock.New()
	defer d.Close()

	rows := sqlmock.NewRows([]string{"id", "category", "resource", "action"}).
		AddRow(1, "auth", "user", "get").
		AddRow(2, "auth", "user", "add")

	mock.ExpectQuery("SELECT id, category, resource, action FROM permissions").WillReturnRows(rows)

	db := SqlDatabase{
		database: d,
		logger:   logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "warn"}),
	}

	permissions, err := db.GetPermissions(nil, context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if len(permissions) != 2 {
		t.Fatal("wrong number of permissions returned")
	}
}

func TestGetPermission(t *testing.T) {
	d, mock, _ := sqlmock.New()
	defer d.Close()

	rows := sqlmock.NewRows([]string{"id", "category", "resource", "action"}).
		AddRow(1, "auth", "user", "get")

	mock.ExpectQuery(`SELECT id, category, resource, action FROM permissions WHERE category = \$1`).
		WithArgs("auth").
		WillReturnRows(rows)

	db := SqlDatabase{
		database: d,
		logger:   logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "warn"}),
	}

	permission, err := db.GetPermission(Filter{
		Key:      "category",
		Operator: "=",
		Value:    "auth",
	}, context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(permission, Permission{ID: 1, Category: "auth", Resource: "user", Action: "get"}) {
		t.Fatal("wrong permission returned")
	}
}
