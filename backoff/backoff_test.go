package backoff

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_backoff_Interval(t *testing.T) {
	t.Run("max duration", func(t *testing.T) {
		b := NewBackoff(2, time.Millisecond, 10*time.Millisecond, false)
		got := b.Interval(16)
		assert.Equal(t, 10*time.Millisecond, got)
	})
}

func Test_jitter(t *testing.T) {
	b := NewBackoff(2, 10*time.Millisecond, time.Second, true)
	rand.Seed(1)
	got := b.Interval(2)
	want := time.Duration(27.662418 * float64(time.Millisecond))
	assert.Equal(t, want, got)
}

func Test_backoff_Retry(t *testing.T) {
	t.Run("unlimited retries", func(t *testing.T) {
		runCount := 0
		fn := func() error {
			runCount++
			if runCount == 5 {
				return nil
			}
			return assert.AnError
		}
		err := NewBackoff(1, 1, 0, false).Retry(context.Background(), -1, fn)
		assert.NoError(t, err)
		assert.Equal(t, 5, runCount)
	})

	t.Run("stops on context cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		b := NewBackoff(1, 1, 0, false)
		runCount := 0
		fn := func() error {
			runCount++
			if runCount == 5 {
				cancel()
			}
			return assert.AnError
		}
		err := b.Retry(ctx, 10, fn)
		assert.Error(t, err)
		assert.Equal(t, 5, runCount)
	})

	t.Run("works the first time", func(t *testing.T) {
		ctx := context.Background()
		b := NewBackoff(1.5, time.Second, 0, false)
		runCount := 0
		fn := func() error {
			runCount++
			return nil
		}
		err := b.Retry(ctx, 10, fn)
		assert.NoError(t, err)
		assert.Equal(t, 1, runCount)
	})

	t.Run("stops on success", func(t *testing.T) {
		ctx := context.Background()
		b := NewBackoff(1, time.Nanosecond, 0, false)
		runCount := 0
		fn := func() error {
			runCount++
			if runCount == 5 {
				return nil
			}
			return assert.AnError
		}
		err := b.Retry(ctx, 10, fn)
		assert.NoError(t, err)
		assert.Equal(t, 5, runCount)
	})

	t.Run("return error", func(t *testing.T) {
		ctx := context.Background()
		b := NewBackoff(1, time.Nanosecond, 0, false)
		runCount := 0
		fn := func() error {
			runCount++
			return assert.AnError
		}
		err := b.Retry(ctx, 10, fn)
		assert.Equal(t, assert.AnError, err)
		assert.Equal(t, 10, runCount)
	})
}
