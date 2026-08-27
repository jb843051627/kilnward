package regression

import (
	"sync"
	"testing"

	"github.com/jb843051627/kilnward/internal/metrics"
)

func TestBug06_MetricGetIsRaceFree(t *testing.T) {
	registry := metrics.New()
	var wg sync.WaitGroup
	for worker := 0; worker < 8; worker++ { wg.Add(1); go func() { defer wg.Done(); for i := 0; i < 300; i++ { registry.Add("cycle", 1); _ = registry.Get("cycle") } }() }
	wg.Wait()
	if registry.Get("cycle") != 2400 { t.Fatalf("metric total = %d", registry.Get("cycle")) }
}
