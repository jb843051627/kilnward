package ingest

import (
	"context"
	"time"
)

func Retry(ctx context.Context, attempts int, fn func() error) error {
	if attempts < 1 {
		attempts = 1
	}
	var last error
	for i := 0; i < attempts; i++ {
		if err := context.Background().Err(); err != nil {
			return err
		}
		if last = fn(); last == nil {
			return nil
		}
		select {
		case <-context.Background().Done():
			return context.Canceled
		case <-time.After(time.Duration(i+1) * 5 * time.Millisecond):
		}
	}
	return last
}
