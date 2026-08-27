package regression

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/store"
)

func TestBug19_MissingKilnMapsToDomainNotFound(t *testing.T) {
	repo, err := store.Open(filepath.Join(t.TempDir(), "case.db")); if err != nil { t.Fatal(err) }; defer repo.Close()
	_, err = repo.GetKiln(context.Background(), "missing-kiln")
	if !errors.Is(err, model.ErrNotFound) { t.Fatalf("error = %v, want model.ErrNotFound", err) }
}
