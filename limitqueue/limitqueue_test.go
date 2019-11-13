package limitqueue

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	t.Run("full", func(t *testing.T) {
		queueLength := 2
		var runCount int64
		var waitToStartJobs, waitForAllJobs, waitForFirstJobStart sync.WaitGroup
		var firstJobOnce sync.Once
		job := func() {
			firstJobOnce.Do(func() {
				waitForFirstJobStart.Done()
			})
			waitToStartJobs.Wait()
			runCount++
			waitForAllJobs.Done()
		}
		waitToStartJobs.Add(1)
		waitForFirstJobStart.Add(1)
		waitForAllJobs.Add(queueLength + 1)
		queue := New(queueLength)
		assert.NoError(t, queue.Enqueue(job))
		waitForFirstJobStart.Wait()
		for i := 0; i < queueLength; i++ {
			assert.NoError(t, queue.Enqueue(job))
		}
		err := queue.Enqueue(job)
		assert.Equal(t, ErrFullQueue, err)
		waitToStartJobs.Done()
		waitForAllJobs.Wait()
		assert.Equal(t, int64(queueLength+1), runCount)
		err = queue.Close(context.Background())
		assert.NoError(t, err)
		err = queue.Enqueue(job)
		assert.Equal(t, ErrClosedQueue, err)
	})

	t.Run("close timeout", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)
		job := func() {
			wg.Wait()
		}
		queue := New(1)
		assert.NoError(t, queue.Enqueue(job))
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		err := queue.Close(ctx)
		assert.Error(t, err)
	})
}
