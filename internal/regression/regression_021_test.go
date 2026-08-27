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

func TestBug21_DuplicateTelemetrySequenceIsRejected(t *testing.T) {
	repo, err := store.Open(filepath.Join(t.TempDir(), "case.db")); if err != nil { t.Fatal(err) }; defer repo.Close(); app := service.NewLab(repo); defer app.Close()
	kiln, err := app.RegisterKiln(context.Background(), model.Kiln{Name: "K1", Location: "east", MaxTempC: 1500, ProbeCount: 2}); if err != nil { t.Fatal(err) }; load, err := app.CreateLoad(context.Background(), model.Load{KilnID: kiln.ID, Label: "L1", Material: model.Material{Code: "clay", Quantity: 1}, Profile: "standard", TargetTempC: 900}); if err != nil { t.Fatal(err) }; cycle, err := app.CycleForLoad(context.Background(), load.ID); if err != nil { t.Fatal(err) }
 	if _, err := repo.LastTelemetry(context.Background(), cycle.ID); !errors.Is(err, model.ErrNotFound) { t.Fatalf("empty telemetry error = %v, want not found", err) }
 	frame := model.TelemetryFrame{ID: "frame-1", KilnID: kiln.ID, LoadID: load.ID, CycleID: cycle.ID, Gateway: "gw-1", Sequence: 1, ReceivedAt: time.Now().UTC(), Samples: []model.ProbeSample{{Sensor: "P1", Temperature: 300, Atmosphere: 50, Power: 20, RecordedAt: time.Now().UTC()}}}
	if _, err := app.IngestTelemetry(context.Background(), frame); err != nil { t.Fatal(err) }; _, err = app.IngestTelemetry(context.Background(), frame); if !errors.Is(err, model.ErrConflict) { t.Fatalf("duplicate telemetry error = %v, want conflict", err) }
}
