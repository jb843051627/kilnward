package regression

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"github.com/jb843051627/kilnward/internal/store"
)

func TestBug30_TransactionStopsOnCanceledContext(t *testing.T) {
	repo, err := store.Open(filepath.Join(t.TempDir(), "case.db")); if err != nil { t.Fatal(err) }; defer repo.Close(); canceled, cancel := context.WithCancel(context.Background()); cancel(); called := false
	err = repo.Transaction(canceled, func(*sql.Tx) error { called = true; return nil })
	if !errors.Is(err, context.Canceled) { t.Fatalf("error = %v, want context.Canceled", err) }; if called { t.Fatal("canceled transaction callback ran") }
}
