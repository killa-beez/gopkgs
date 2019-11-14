package pool

import (
	"context"
	"sync"
)

//WorkUnit describes the interface for a work unit
type WorkUnit interface {
	Perform(ctx context.Context)
}

//Pool describes the pool of workers
type Pool struct {
	numWorkers int
	ch         chan WorkUnit
	wg         sync.WaitGroup
}

//New returns a instantiated pool
func New(queueSize, numWorkers int) Pool {
	return Pool{
		numWorkers: numWorkers,
		ch:         make(chan WorkUnit, queueSize),
	}
}

type workUnit func(context.Context)

func (p workUnit) Perform(ctx context.Context) {
	p(ctx)
}

//NewWorkUnit returns a new WorkUnit that will execute fn
func NewWorkUnit(fn func(context.Context)) WorkUnit {
	return workUnit(fn)
}

//Start starts the work
func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.numWorkers; i++ {
		p.wg.Add(1)

		go func() {
			defer p.wg.Done()

			for unit := range p.ch {
				if ctx.Err() != nil {
					break
				}
				unit.Perform(ctx)
			}
		}()
	}
}

//Add adds a work unit
func (p *Pool) Add(obj WorkUnit) {
	p.ch <- obj
}

//Wait closes the channel for adding and waits for the waitgroup
func (p *Pool) Wait() {
	close(p.ch)
	p.wg.Wait()
}
