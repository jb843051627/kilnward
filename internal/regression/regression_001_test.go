package regression

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"

	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/store"
)

func TestBug01_TransactionReturnsCallerError(t *testing.T) {
	repo, err := store.Open(filepath.Join(t.TempDir(), "case.db"))
	if err != nil { t.Fatal(err) }
	defer repo.Close()
	sentinel := errors.New("stage gate rejected")
	err = repo.Transaction(context.Background(), func(*sql.Tx) error { return sentinel })
	if !errors.Is(err, sentinel) { t.Fatalf("transaction error = %v, want caller error", err) }
	profile := model.Profile{ID: "p-1", Name: "profile", Material: "clay", Author: "operator", Steps: []model.ProfileStep{{Sequence: 0, Name: "预热", TargetTempC: 300, RampPerMinute: 5}, {Sequence: 1, Name: "升温", TargetTempC: 900, RampPerMinute: 5}}}
	if err := repo.CreateProfile(context.Background(), profile); err != nil { t.Fatal(err) }
	if err := repo.CreateProfile(context.Background(), profile); err == nil { t.Fatal("duplicate profile write was accepted") }
}
