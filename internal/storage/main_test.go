package storage

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	db          *sqlx.DB
	storageRepo *Storage
)

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

func runTests(m *testing.M) int {
	var err error
	ctx := context.Background()

	testDatabase, err := newTestDatabase(ctx)
	if err != nil {
		log.Fatalf("error creating test database: %v", err)
	}
	defer testDatabase.Terminate(ctx)

	connectionString, err := testDatabase.ConnectionString(context.Background(), "sslmode=disable")
	if err != nil {
		log.Fatalf("error getting connection string: %v", err)
	}

	db, err = sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	storageRepo, err = New(db)
	if err != nil {
		log.Fatal(err)
	}

	return m.Run()
}

func newTestDatabase(ctx context.Context) (*postgres.PostgresContainer, error) {
	pgContainer, err := postgres.Run(ctx, "postgres:15.3-alpine",
		postgres.WithInitScripts("./../../init.sql"),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	return pgContainer, nil
}

func prepareDB(queries ...string) error {
	for _, q := range queries {
		_, err := db.Exec(q)
		if err != nil {
			return err
		}
	}
	return nil
}

func clearTableUsers() error {
	_, err := db.Exec("delete from users")
	return err
}
