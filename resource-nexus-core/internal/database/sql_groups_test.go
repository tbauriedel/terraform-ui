package database

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tbauriedel/resource-nexus-core/internal/config"
	"github.com/tbauriedel/resource-nexus-core/internal/logging"
)

func TestGetGroups(t *testing.T) {
	d, mock, _ := sqlmock.New()
	defer d.Close()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "dummy").
		AddRow(2, "dummy2")

	mock.ExpectQuery("SELECT id, name FROM groups").WillReturnRows(rows)

	db := SqlDatabase{
		database: d,
		logger:   logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "warn"}),
	}

	groups, err := db.GetGroups(nil, context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if len(groups) != 2 {
		t.Fatal("expected 2 groups")
	}

	if groups[0].Name != "dummy" || groups[1].Name != "dummy2" {
		t.Fatal("wrong groups returned")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("there were unfulfilled expectations: %v", err)
	}
}

func TestGetGroup(t *testing.T) {
	d, mock, _ := sqlmock.New()
	defer d.Close()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "dummy")

	db := SqlDatabase{
		database: d,
		logger:   logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "warn"}),
	}

	mock.ExpectQuery(`SELECT id, name FROM groups WHERE name = \$1`).WithArgs("dummy").WillReturnRows(rows)

	ctx := context.TODO()
	group, err := db.GetGroup(
		Filter{
			Key:      "name",
			Operator: "=",
			Value:    "dummy",
		}, ctx,
	)
	if err != nil {
		t.Fatal(err)
	}

	if group.Name != "dummy" {
		t.Fatal("wrong group returned")
	}
}

func TestInsertGroup(t *testing.T) {
	d, mock, _ := sqlmock.New()
	defer d.Close()

	group := Group{
		Name: "dummy",
	}

	mock.ExpectExec(`INSERT INTO groups \(name\) VALUES \(\$1\)`).
		WithArgs(group.Name).
		WillReturnResult(sqlmock.NewResult(1, 1))

	db := SqlDatabase{
		database: d,
		logger:   logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "warn"}),
	}

	ctx := context.TODO()
	res, err := db.InsertGroup(ctx, group)
	if err != nil {
		t.Fatal(err)
	}

	if rows, _ := res.RowsAffected(); rows != 1 {
		t.Fatal("wrong number of rows affected")
	}
}
