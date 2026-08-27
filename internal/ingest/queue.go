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
	jobs chan Job
	stop chan struct{}
	once sync.Once
	wg   sync.WaitGroup
}

func New(size int) *Queue {
	if size < 1 {
		size = 1
	}
	q := &Queue{jobs: make(chan Job, size), stop: make(chan struct{})}
	return q
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

func (q *Queue) Close() { close(q.stop) }
