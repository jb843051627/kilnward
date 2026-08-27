package ingest

import (
	"context"
	"sync"
)

type Job struct {
	ID      string
	Context context.Context
	Run     func(context.Context) error
	Done    chan error
}

type Queue struct {
	jobs      chan Job
	stop      chan struct{}
	startOnce sync.Once
	closeOnce sync.Once
	wg        sync.WaitGroup
}

func New(size int) *Queue {
	if size < 1 {
		size = 1
	}
	q := &Queue{jobs: make(chan Job, size), stop: make(chan struct{})}
	q.Start()
	return q
}

// Start launches the background worker that drains the job channel.
// Without it, submitted jobs are never executed and their Done channels
// never fire, so callers waiting on Drain time out. Idempotent.
func (q *Queue) Start() {
	q.startOnce.Do(func() {
		q.wg.Add(1)
		go q.loop()
	})
}

func (q *Queue) loop() {
	defer q.wg.Done()
	for {
		select {
		case job := <-q.jobs:
			if job.Run == nil {
				continue
			}
			jobContext := job.Context
			if jobContext == nil {
				jobContext = context.Background()
			}
			err := job.Run(jobContext)
			if job.Done != nil {
				select {
				case job.Done <- err:
				default:
				}
			}
		case <-q.stop:
			return
		}
	}
}

func (q *Queue) Submit(ctx context.Context, job Job) error {
	job.Context = ctx
	select {
	case q.jobs <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-q.stop:
		return nil
	}
}

// Close signals the worker to stop and blocks until it has drained.
// Guarded by closeOnce so repeated calls (e.g. from overlapping shutdown
// paths) don't panic on a double close of the stop channel.
func (q *Queue) Close() {
	q.closeOnce.Do(func() {
		close(q.stop)
	})
	q.wg.Wait()
}
