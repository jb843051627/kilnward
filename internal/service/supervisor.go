package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kilnward/internal/ingest"
	"github.com/jb843051627/kilnward/internal/model"
)

func (a *App) QueueReadingCheck(ctx context.Context, reading model.Reading) (<-chan error, error) {
	if err := a.requireContext(ctx); err != nil {
		return nil, err
	}
	done := make(chan error, 1)
	job := ingest.Job{ID: reading.ID, Done: done, Run: func(jobCtx context.Context) error {
		_, err := a.RecordReading(jobCtx, reading)
		if err != nil {
			return fmt.Errorf("reading job %s: %w", reading.ID, err)
		}
		return nil
	}}
	if err := a.queue.Submit(ctx, job); err != nil {
		return nil, err
	}
	return done, nil
}

func (a *App) DrainReadingCheck(ctx context.Context, done <-chan error) error {
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
