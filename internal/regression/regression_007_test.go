package regression

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/jb843051627/kilnward/internal/ingest"
)

func TestBug07_QueueStartsWorkerAndClosesOnce(t *testing.T) {
	q := ingest.New(1)
	done := make(chan error, 1)
	if err := q.Submit(context.Background(), ingest.Job{Done: done, Run: func(context.Context) error { return nil }}); err != nil { t.Fatal(err) }
	select { case err := <-done: if err != nil { t.Fatal(err) }; case <-time.After(time.Second): t.Fatal("worker did not run") }
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ { wg.Add(1); go func() { defer wg.Done(); q.Close() }() }
	wg.Wait()
}
