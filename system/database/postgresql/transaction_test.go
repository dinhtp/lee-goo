package postgresql

import (
	"errors"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// openClosedDB returns a sqlx.DB that has been closed — calling Beginx on it returns an error.
// This avoids requiring a real database for the begin-error path.
func openClosedDB(t *testing.T) *sqlx.DB {
	t.Helper()

	// sqlx.Open is lazy; it does not dial the server, so a bad DSN still succeeds here.
	db, err := sqlx.Open("pgx", "host=127.0.0.1 port=1 user=x password=x dbname=x sslmode=disable")
	if err != nil {
		t.Fatalf("sqlx.Open: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("db.Close: %v", err)
	}
	return db
}

func TestTransactReturnsBeginError(t *testing.T) {
	db := openClosedDB(t)

	err := Transact(db, func(_ *sqlx.Tx) error {
		t.Fatal("callback must not run when Begin fails")
		return nil
	})
	if err == nil {
		t.Fatal("expected begin error, got nil")
	}
}

func TestTransactRollsBackOnCallbackError(t *testing.T) {
	if testing.Short() {
		t.Skip("requires a running PostgreSQL instance")
	}

	dsn := "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=lee_goo sslmode=disable"
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS transact_test (id SERIAL PRIMARY KEY, name TEXT)`); err != nil {
		t.Fatalf("create table: %v", err)
	}
	defer db.Exec(`DROP TABLE IF EXISTS transact_test`)

	callbackErr := errors.New("intentional failure")
	err = Transact(db, func(tx *sqlx.Tx) error {
		if _, err := tx.Exec(`INSERT INTO transact_test (name) VALUES ($1)`, "rollback"); err != nil {
			return err
		}
		return callbackErr
	})
	if !errors.Is(err, callbackErr) {
		t.Fatalf("expected %v, got %v", callbackErr, err)
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM transact_test WHERE name = $1`, "rollback").Scan(&count); err != nil {
		t.Fatalf("count query: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected rollback to leave 0 rows, got %d", count)
	}
}
