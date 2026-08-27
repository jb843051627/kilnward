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

func TestBug18_CriticalIncidentReturnsQuarantineError(t *testing.T) {
	repo, err := store.Open(filepath.Join(t.TempDir(), "case.db")); if err != nil { t.Fatal(err) }; defer repo.Close(); app := service.NewLab(repo); defer app.Close()
	kiln, err := app.RegisterKiln(context.Background(), model.Kiln{Name: "offline", Location: "west", State: model.KilnOffline, MaxTempC: 1400, ProbeCount: 2}); if err != nil { t.Fatal(err) }
	_, err = app.OpenIncident(context.Background(), model.Incident{KilnID: kiln.ID, Code: "overheat", Severity: model.SeverityCritical, Detail: "temperature exceeded"})
	if !errors.Is(err, model.ErrInvalidState) { t.Fatalf("error = %v, want invalid state", err) }
}
