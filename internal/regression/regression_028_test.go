package regression

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/service"
	"github.com/jb843051627/kilnward/internal/store"
)

func TestBug28_CanceledTelemetryDoesNotPersist(t *testing.T) {
	repo, err := store.Open(filepath.Join(t.TempDir(), "case.db")); if err != nil { t.Fatal(err) }; defer repo.Close(); app := service.NewLab(repo); defer app.Close(); ctx := context.Background(); kiln, err := app.RegisterKiln(ctx, model.Kiln{Name: "K1", Location: "east", MaxTempC: 1500, ProbeCount: 2}); if err != nil { t.Fatal(err) }; load, err := app.CreateLoad(ctx, model.Load{KilnID: kiln.ID, Label: "L1", Material: model.Material{Code: "clay", Quantity: 1}, Profile: "standard", TargetTempC: 900}); if err != nil { t.Fatal(err) }; cycle, err := app.CycleForLoad(ctx, load.ID); if err != nil { t.Fatal(err) }; canceled, cancel := context.WithCancel(ctx); cancel()
	_, err = app.IngestTelemetry(canceled, model.TelemetryFrame{ID: "f1", KilnID: kiln.ID, LoadID: load.ID, CycleID: cycle.ID, Gateway: "gw", Sequence: 1, ReceivedAt: time.Now().UTC(), Samples: []model.ProbeSample{{Sensor: "P1", Temperature: 300, Atmosphere: 50, Power: 10, RecordedAt: time.Now().UTC()}}})
	if !errors.Is(err, context.Canceled) { t.Fatalf("error = %v, want context.Canceled", err) }
}
