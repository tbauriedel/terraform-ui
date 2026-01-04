package database

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tbauriedel/resource-nexus-core/internal/config"
	"github.com/tbauriedel/resource-nexus-core/internal/logging"
)

func TestGetGroupPermissions(t *testing.T) {
	d, mock, _ := sqlmock.New()
	defer d.Close()

	rows := sqlmock.NewRows([]string{"group_id", "permission_id"}).
		AddRow(1, 1).
		AddRow(2, 2)

	mock.ExpectQuery("SELECT group_id, permission_id FROM group_permissions").
		WillReturnRows(rows)

	db := SqlDatabase{
		database: d,
		logger:   logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "warn"}),
	}

	groupPermissions, err := db.GetGroupPermissions(nil, context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if len(groupPermissions) != 2 {
		t.Fatal("wrong number of group permissions returned")
	}
}

func TestGetGroupPermission(t *testing.T) {
	d, mock, _ := sqlmock.New()
	defer d.Close()

	rows := sqlmock.NewRows([]string{"group_id", "permission_id"}).
		AddRow(1, 4)

	mock.ExpectQuery(`SELECT group_id, permission_id FROM group_permissions WHERE \(group_id = \$1 AND permission_id = \$2\)`).
		WithArgs(1, 4).
		WillReturnRows(rows)

	db := SqlDatabase{
		database: d,
		logger:   logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "warn"}),
	}

	groupPermission, err := db.GetGroupPermission(LogicalFilter{
		Operator: "AND",
		Filters: []FilterExpr{
			Filter{
				Key:      "group_id",
				Operator: "=",
				Value:    1,
			},
			Filter{
				Key:      "permission_id",
				Operator: "=",
				Value:    4,
			},
		},
	}, context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if groupPermission.GroupID != 1 || groupPermission.PermissionID != 4 {
		t.Fatal("wrong group permission returned")
	}
}

func TestInsertGroupPermission(t *testing.T) {
	d, mock, _ := sqlmock.New()
	defer d.Close()

	groupPermission := GroupPermissionReference{
		GroupID:      1,
		PermissionID: 4,
	}

	mock.ExpectExec(`INSERT INTO group_permissions \(group_id\, permission_id\) VALUES \(\$1, \$2\)`).
		WithArgs(1, 4).
		WillReturnResult(sqlmock.NewResult(1, 1))

	db := SqlDatabase{
		database: d,
		logger:   logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "warn"}),
	}

	_, err := db.InsertGroupPermission(context.TODO(), groupPermission)
	if err != nil {
		t.Fatal(err)
	}
}
