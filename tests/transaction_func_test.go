package tests

import (
	"context"
	"database/sql"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"monyLonger/internal/storage"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestTransactionStorage_GetByVaultID(t *testing.T) {
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx, "postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer pgContainer.Terminate(ctx)

	dsn, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := storage.RunMigrations(db, migrationsPath()); err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO vaults (id, name, limit_amount, left_amount, expire) 
        VALUES (1, 'Продукты', 8000, 8000, NOW())`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO transactions (name, amount, is_income, vault_id, note, date) VALUES
                                                                 ('La Pore',1000,false,1,' ','2026-06-15 00:00:00.000000'),
                                                                 ('Li Fo',1000,false,1,' ','2026-05-15 00:00:00.000000')`)
	if err != nil {
		t.Fatal(err)
	}

	storage := storage.NewTransactionStorage(db)

	result, err := storage.GetForPeriod(time.Date(2026, time.Month(6), 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 transaction, got %d", len(result))
	}
}

func migrationsPath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return "file://" + filepath.Join(dir, "..", "migrations")
}
