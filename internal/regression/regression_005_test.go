package regression

import (
	"sync"
	"testing"

	"github.com/jb843051627/kilnward/internal/metrics"
)

func TestBug05_MetricUpdatesAreRaceFree(t *testing.T) {
	registry := metrics.New()
	var wg sync.WaitGroup
	for worker := 0; worker < 6; worker++ { wg.Add(1); go func() { defer wg.Done(); for i := 0; i < 200; i++ { registry.Add("temperature", 1); _ = registry.Snapshot() } }() }
	wg.Wait()
	if registry.Get("temperature") != 1200 { t.Fatalf("metric total = %d", registry.Get("temperature")) }
}
