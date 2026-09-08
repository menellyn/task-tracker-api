package task

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/menellyn/task-tracker-api/internal/database"
	postgres_container "github.com/testcontainers/testcontainers-go/modules/postgres"
)

var (
	testDB        *sql.DB
	testContainer *postgres_container.PostgresContainer
	testGormDB    *gorm.DB
)

func TestMain(m *testing.M) {

	ctx := context.Background()

	var err error

	testContainer, err = postgres_container.Run(
		ctx,
		"postgres:16-alpine",
		postgres_container.WithDatabase("task_tracker_test"),
		postgres_container.WithUsername("postgres"),
		postgres_container.WithPassword("postgres"),
		postgres_container.BasicWaitStrategies(),
		postgres_container.WithSQLDriver("pgx"),
	)
	if err != nil {
		log.Fatal(err)
	}

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

	if err := testDB.Ping(); err != nil {
		log.Fatal(err)
	}

	testGormDB, err = gorm.Open(
		gorm_postgres.New(gorm_postgres.Config{
			Conn: testDB,
		}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.RunMigrations(testDB); err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	if err := testDB.Close(); err != nil {
		log.Printf("failed to close test DB: %v", err)
	}

	if err := testContainer.Terminate(ctx); err != nil {
		log.Printf("failed to terminate test container: %v", err)
	}

	os.Exit(code)
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

func newTestGormTx(t *testing.T) *gorm.DB {
	t.Helper()

	tx := testGormDB.Begin()
	if tx.Error != nil {
		t.Fatal(tx.Error)
	}

	t.Cleanup(func() {
		if err := tx.Rollback().Error; err != nil {
			t.Logf("failed to rollback transaction: %v", err)
		}
	})

	return tx
}
