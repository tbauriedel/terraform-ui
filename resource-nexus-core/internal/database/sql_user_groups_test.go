package database

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tbauriedel/resource-nexus-core/internal/config"
	"github.com/tbauriedel/resource-nexus-core/internal/logging"
)

func TestGetUserGroupReferences(t *testing.T) {
	d, mock, _ := sqlmock.New()
	defer d.Close()

	rows := sqlmock.NewRows([]string{"user_id", "group_id"}).
		AddRow(7, 14).
		AddRow(8, 14)

	mock.ExpectQuery(`SELECT user_id, group_id FROM user_groups`).WillReturnRows(rows)

	db := SqlDatabase{
		database: d,
		logger:   logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "warn"}),
	}

	userGroupRefs, err := db.GetUserGroupReferences(nil, context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if len(userGroupRefs) != 2 {
		t.Fatal("wrong number of user group references returned")
	}
}

func TestGetUserGroupReference(t *testing.T) {
	d, mock, _ := sqlmock.New()
	defer d.Close()

	rows := sqlmock.NewRows([]string{"user_id", "group_id"}).
		AddRow(7, 14)

	mock.ExpectQuery(`SELECT user_id, group_id FROM user_groups`).WillReturnRows(rows)

	db := SqlDatabase{
		database: d,
		logger:   logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "warn"}),
	}

	ctx := context.TODO()
	userGroupRef, err := db.GetUserGroupReference(nil, ctx)
	if err != nil {
		t.Fatal(err)
	}

	if userGroupRef.UserID != 7 || userGroupRef.GroupID != 14 {
		t.Fatal("wrong user group reference returned")
	}
}

func TestInsertUserGroupReference(t *testing.T) {
	d, mock, _ := sqlmock.New()
	defer d.Close()

	usergroup := UserGroupReference{
		UserID:  7,
		GroupID: 14,
	}

	mock.ExpectExec(`INSERT INTO user_groups \(user_id\, group_id\) VALUES \(\$1, \$2\)`).
		WithArgs(7, 14).
		WillReturnResult(sqlmock.NewResult(1, 1))

	db := SqlDatabase{
		database: d,
		logger:   logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "warn"}),
	}

	ctx := context.TODO()
	res, err := db.InsertUserGroupReference(ctx, usergroup)
	if err != nil {
		t.Fatal(err)
	}

	if rows, _ := res.RowsAffected(); rows != 1 {
		t.Fatal("wrong number of rows affected")
	}
}
