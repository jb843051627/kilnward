package regression

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"github.com/jb843051627/kilnward/internal/service"
	"github.com/jb843051627/kilnward/internal/store"
)

func TestBug20_ReadingWindowHonorsCanceledContext(t *testing.T) {
	repo, err := store.Open(filepath.Join(t.TempDir(), "case.db")); if err != nil { t.Fatal(err) }; defer repo.Close(); app := service.NewLab(repo); defer app.Close()
	ctx, cancel := context.WithCancel(context.Background()); cancel()
	_, err = app.ReadingWindow(ctx, "cycle-missing")
	if !errors.Is(err, context.Canceled) { t.Fatalf("error = %v, want context.Canceled", err) }
}
