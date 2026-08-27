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

func TestBug25_RecordReadingRejectsInvalidTemperature(t *testing.T) {
	repo, err := store.Open(filepath.Join(t.TempDir(), "case.db")); if err != nil { t.Fatal(err) }; defer repo.Close(); app := service.NewLab(repo); defer app.Close(); ctx := context.Background(); kiln, err := app.RegisterKiln(ctx, model.Kiln{Name: "K1", Location: "east", MaxTempC: 1500, ProbeCount: 2}); if err != nil { t.Fatal(err) }; load, err := app.CreateLoad(ctx, model.Load{KilnID: kiln.ID, Label: "L1", Material: model.Material{Code: "clay", Quantity: 1}, Profile: "standard", TargetTempC: 900}); if err != nil { t.Fatal(err) }; cycle, err := app.CycleForLoad(ctx, load.ID); if err != nil { t.Fatal(err) }
	_, err = app.RecordReading(ctx, model.Reading{KilnID: kiln.ID, LoadID: "other-load", CycleID: cycle.ID, Sensor: "P1", Temperature: 300, Atmosphere: 50, Power: 10})
	if !errors.Is(err, model.ErrConflict) { t.Fatalf("error = %v, want conflict", err) }
	if err := app.RecordReadings(ctx, []model.Reading{{KilnID: kiln.ID, LoadID: load.ID, CycleID: cycle.ID, Sensor: "P2", Temperature: 3000, Atmosphere: 50, Power: 10}}); err == nil { t.Fatal("invalid batch reading was accepted") }
}
