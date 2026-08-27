package regression

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/jb843051627/kilnward/internal/service"
	"github.com/jb843051627/kilnward/internal/store"
)

func TestBug02_HealthHonorsCanceledContext(t *testing.T) {
	repo, err := store.Open(filepath.Join(t.TempDir(), "case.db"))
	if err != nil { t.Fatal(err) }
	defer repo.Close()
	app := service.NewLab(repo)
	defer app.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := app.Health(ctx); !errors.Is(err, context.Canceled) { t.Fatalf("health error = %v, want context.Canceled", err) }
}
