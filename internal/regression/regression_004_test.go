package regression

import (
	"context"
	"errors"
	"testing"

	"github.com/jb843051627/kilnward/internal/ingest"
)

func TestBug04_RetryStopsWhenContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0
	sentinel := errors.New("gateway unavailable")
	err := ingest.Retry(ctx, 5, func() error { calls++; return sentinel })
	if !errors.Is(err, context.Canceled) { t.Fatalf("retry error = %v, want context.Canceled", err) }
	if calls != 0 { t.Fatalf("retry called operation %d times after cancellation", calls) }
}
