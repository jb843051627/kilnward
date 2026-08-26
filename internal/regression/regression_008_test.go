package regression

import (
	"path/filepath"
	"testing"

	"github.com/jb843051627/kilnward/internal/metrics"
	"github.com/jb843051627/kilnward/internal/service"
	"github.com/jb843051627/kilnward/internal/store"
)

func TestBug08_MetricSnapshotDoesNotAliasRegistry(t *testing.T) {
	registry := metrics.New(); registry.Add("temperature", 3)
	snapshot := registry.Snapshot(); snapshot["temperature"] = 99
	if registry.Get("temperature") != 3 { t.Fatalf("registry was changed through snapshot: %v", registry.Get("temperature")) }
	repo, err := store.Open(filepath.Join(t.TempDir(), "case.db")); if err != nil { t.Fatal(err) }; defer repo.Close()
	app := service.NewLab(repo); defer app.Close()
	returned := app.MetricsSnapshot(); returned["temporary"] = 1
	if app.MetricsSnapshot()["temporary"] != 0 { t.Fatal("summary snapshot changed internal metrics") }
}
