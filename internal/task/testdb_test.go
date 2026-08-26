package task

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	testDB        *sql.DB
	testContainer *postgres.PostgresContainer
)

func TestMain(m *testing.M) {

	ctx := context.Background()

	var err error

	testContainer, err = postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("task_tracker_test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	if err != nil {
		log.Fatal(err)
	}

	defer testContainer.Terminate(ctx)

	connString, err := testContainer.ConnectionString(
		ctx,
		"sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}

	testDB, err = sql.Open("pgx", connString)
	if err != nil {
		log.Fatal(err)
	}

	defer testDB.Close()

	if err := testDB.Ping(); err != nil {
		log.Fatal(err)
	}

	createSchema(testDB)

	code := m.Run()

	os.Exit(code)
}

func createSchema(db *sql.DB) {

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			title TEXT NOT NULL,
			done BOOLEAN NOT NULL DEFAULT FALSE
		)
	`)

	if err != nil {
		log.Fatal(err)
	}
}

func newTestTx(t *testing.T) *sql.Tx {
	t.Helper()

	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			t.Logf("failed to rollback transaction: %v", err)
		}
	})
	return tx
}
