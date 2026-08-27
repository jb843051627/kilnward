package regression

import (
	"context"
	"path/filepath"
	"testing"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/service"
	"github.com/jb843051627/kilnward/internal/store"
)

func TestBug22_BatchReadingsRollbackTogether(t *testing.T) {
	repo, err := store.Open(filepath.Join(t.TempDir(), "case.db")); if err != nil { t.Fatal(err) }; defer repo.Close(); app := service.NewLab(repo); defer app.Close(); ctx := context.Background()
	kiln, err := app.RegisterKiln(ctx, model.Kiln{Name: "K1", Location: "east", MaxTempC: 1500, ProbeCount: 2}); if err != nil { t.Fatal(err) }; load, err := app.CreateLoad(ctx, model.Load{KilnID: kiln.ID, Label: "L1", Material: model.Material{Code: "clay", Quantity: 1}, Profile: "standard", TargetTempC: 900}); if err != nil { t.Fatal(err) }; cycle, err := app.CycleForLoad(ctx, load.ID); if err != nil { t.Fatal(err) }
	items := []model.Reading{{ID: "r-1", KilnID: kiln.ID, LoadID: load.ID, CycleID: cycle.ID, Sensor: "P1", Temperature: 300, Atmosphere: 50, Power: 20}, {ID: "r-2", KilnID: kiln.ID, LoadID: load.ID, CycleID: cycle.ID, Sensor: "P2", Temperature: 3000, Atmosphere: 50, Power: 20}}
	if err := app.RecordReadings(ctx, items); err == nil { t.Fatal("invalid batch unexpectedly succeeded") }; stored, err := app.ListReadings(ctx, cycle.ID, 100); if err != nil { t.Fatal(err) }; if len(stored) != 0 { t.Fatalf("partial batch persisted %d readings", len(stored)) }
}
