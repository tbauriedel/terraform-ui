package database

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/tbauriedel/resource-nexus-core/internal/config"
	"github.com/tbauriedel/resource-nexus-core/internal/logging"
)

func TestSelect(t *testing.T) {
	d, mock, _ := sqlmock.New()
	defer d.Close()

	// using this query for the test case.
	query := fmt.Sprintf(`SELECT * FROM users`)

	rows := sqlmock.NewRows([]string{"id", "name", "password_hash", "is_admin"}).
		AddRow(1, "dummy", "foobar", false).
		AddRow(2, "dummy2", "foobar2", true)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

	db := SqlDatabase{
		database: d,
		logger:   logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "warn"}),
	}

	ctx := context.TODO()
	_, closeRows, err := db.Select(query, nil, ctx)
	if err != nil {
		t.Fatal(err)
	}

	defer closeRows()
}
