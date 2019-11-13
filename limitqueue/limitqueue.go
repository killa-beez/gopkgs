package limitqueue

import (
	"context"
	"fmt"
	"sync"
)

// queue errors
var (
	ErrFullQueue   = fmt.Errorf("queue is full")
	ErrClosedQueue = fmt.Errorf("queue is closed")
)

//Queue is a job queue
type Queue interface {
	Enqueue(job func()) error
	Close(ctx context.Context) error
}

//New returns a new Queue
func New(maxSize int) Queue {
	q := &queue{
		jobs: make(chan func(), maxSize),
		done: make(chan error),
	}
	go func() {
		for job := range q.jobs {
			job()
		}
		close(q.done)
	}()
	return q
}

type queue struct {
	jobs   chan func()
	done   chan error
	mux    sync.RWMutex
	closed bool
}

func (q *queue) Enqueue(job func()) error {
	q.mux.RLock()
	defer q.mux.RUnlock()
	if q.closed {
		return ErrClosedQueue
	}
	select {
	case q.jobs <- job:
		return nil
	default:
		return ErrFullQueue
	}
}

func (q *queue) Close(ctx context.Context) error {
	q.mux.Lock()
	q.closed = true
	q.mux.Unlock()
	var wg sync.WaitGroup
	wg.Add(1)
	var err error
	go func() {
		if q.done != nil {
			select {
			case <-q.done:
			case <-ctx.Done():
				err = ctx.Err()
			}
		}
		wg.Done()
	}()
	if q.jobs != nil {
		close(q.jobs)
	}
	wg.Wait()
	return err
}
