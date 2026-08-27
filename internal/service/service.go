package service

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/jb843051627/kilnward/internal/clock"
	"github.com/jb843051627/kilnward/internal/ingest"
	"github.com/jb843051627/kilnward/internal/metrics"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/store"
)

type App struct {
	repo     *store.Store
	clock    clock.Clock
	queue    *ingest.Queue
	metrics  *metrics.Registry
	sequence atomic.Int64
}

func NewLab(repo *store.Store) *App {
	app := &App{repo: repo, clock: clock.System{}, queue: ingest.New(64), metrics: metrics.New()}
	app.sequence.Store(time.Now().UnixNano())
	return app
}

func (a *App) Close() {
	if a.queue != nil {
		a.queue.Close()
	}
}

func (a *App) now() time.Time { return a.clock.Now() }

func (a *App) newID(prefix string) string { return model.MakeID(prefix, a.now(), a.sequence.Add(1)) }

func (a *App) audit(ctx context.Context, subjectType, subjectID, action, detail string) error {
	a.metrics.Add("audit.total", 1)
	return a.repo.RecordAudit(ctx, model.Audit{SubjectType: subjectType, SubjectID: subjectID, Action: action, Detail: detail, CreatedAt: a.now()})
}

func (a *App) requireContext(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("%w: nil context", model.ErrValidation)
	}
	return ctx.Err()
}

func (a *App) Metrics() map[string]int64 { return a.metrics.Snapshot() }

func (a *App) Health(ctx context.Context) error {
	if err := a.requireContext(ctx); err != nil {
		return err
	}
	return a.repo.Ping(context.Background())
}
