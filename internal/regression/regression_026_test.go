package regression

 import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/service"
	"github.com/jb843051627/kilnward/internal/store"
)

func TestBug26_ReadingIncidentReturnsMissingStageError(t *testing.T) {
	repo, err := store.Open(filepath.Join(t.TempDir(), "case.db")); if err != nil { t.Fatal(err) }; defer repo.Close(); app := service.NewLab(repo); defer app.Close(); ctx := context.Background(); kiln, err := app.RegisterKiln(ctx, model.Kiln{Name: "K1", Location: "east", MaxTempC: 1500, ProbeCount: 2}); if err != nil { t.Fatal(err) }
	_, err = app.EvaluateReadingIncident(ctx, kiln.ID, "load-missing", model.Reading{CycleID: "cycle-missing", Temperature: 900})
	if !errors.Is(err, model.ErrNotFound) { t.Fatalf("error = %v, want not found", err) }
}
