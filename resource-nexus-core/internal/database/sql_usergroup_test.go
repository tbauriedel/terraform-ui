package database

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tbauriedel/resource-nexus-core/internal/config"
	"github.com/tbauriedel/resource-nexus-core/internal/logging"
)

func TestInsertUserGroupReference(t *testing.T) {
	d, mock, _ := sqlmock.New()
	defer d.Close()

	usergroup := UserGroup{
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
