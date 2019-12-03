package retry

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		ctx := context.Background()
		config := Config{}
		runCount := 0
		wantCount := 1
		err := Retry(ctx, func() error {
			runCount++
			return nil
		}, config)
		assert.NoError(t, err)
		assert.Equal(t, wantCount, runCount)
	})

	t.Run("non-retryable error", func(t *testing.T) {
		ctx := context.Background()
		config := Config{
			Retryable: func(lastErr error) bool {
				t.Helper()
				assert.Equal(t, assert.AnError, lastErr)
				return false
			},
		}
		runCount := 0
		wantCount := 1
		err := Retry(ctx, func() error {
			runCount++
			return assert.AnError
		}, config)
		assert.Equal(t, assert.AnError, err)
		assert.Equal(t, wantCount, runCount)
	})

	t.Run("success on retry", func(t *testing.T) {
		ctx := context.Background()
		config := Config{}
		runCount := 0
		wantCount := 2
		err := Retry(ctx, func() error {
			runCount++
			if runCount > 1 {
				return nil
			}
			return assert.AnError
		}, config)
		assert.NoError(t, err)
		assert.Equal(t, wantCount, runCount)
	})

	t.Run("runs out of attempts", func(t *testing.T) {
		ctx := context.Background()
		config := Config{}
		runCount := 0
		wantCount := 3
		err := Retry(ctx, func() error {
			runCount++
			return assert.AnError
		}, config)
		assert.Equal(t, assert.AnError, err)
		assert.Equal(t, wantCount, runCount)
	})

	t.Run("sleeps", func(t *testing.T) {
		ctx := context.Background()
		var wantSleep time.Duration
		config := Config{
			Backoff: func(retryCount int, _ error) time.Duration {
				sleepTime := time.Duration(retryCount) * time.Millisecond
				wantSleep += sleepTime
				return sleepTime
			},
			MaxAttempts: 10,
		}
		runCount := 0
		wantCount := 5
		startTime := time.Now()
		err := Retry(ctx, func() error {
			runCount++
			if runCount == wantCount {
				return nil
			}
			return assert.AnError
		}, config)
		elapsed := time.Since(startTime)
		assert.NoError(t, err)
		assert.Equal(t, wantCount, runCount)
		assert.GreaterOrEqual(t, elapsed.Milliseconds(), wantSleep.Milliseconds())
		assert.LessOrEqual(t, elapsed.Milliseconds(), 2*wantSleep.Milliseconds())
	})

	t.Run("cancel ctx", func(t *testing.T) {
		ctx := context.Background()
		wantSleep := 10 * time.Millisecond
		config := Config{
			Backoff: func(retryCount int, _ error) time.Duration {
				return time.Second
			},
			MaxAttempts: 10,
		}
		runCount := 0
		wantCount := 1
		startTime := time.Now()
		ctx, cancel := context.WithTimeout(ctx, wantSleep)
		defer cancel()
		err := Retry(ctx, func() error {
			runCount++
			return assert.AnError
		}, config)
		elapsed := time.Since(startTime)
		assert.Equal(t, assert.AnError, err)
		assert.Equal(t, wantCount, runCount)
		assert.GreaterOrEqual(t, elapsed.Milliseconds(), wantSleep.Milliseconds())
		assert.LessOrEqual(t, elapsed.Milliseconds(), 2*wantSleep.Milliseconds())
	})
}

func Test_sleep(t *testing.T) {
	t.Run("sleeps", func(t *testing.T) {
		ctx := context.Background()
		sleepyTime := 2 * time.Millisecond
		startTime := time.Now()
		sleep(ctx, sleepyTime)
		took := time.Since(startTime)
		assert.True(t, took > sleepyTime)
	})

	t.Run("canceled context wakes it up", func(t *testing.T) {
		startTime := time.Now()
		sleepyTime := 200 * time.Millisecond
		cancelTime := 2 * time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), cancelTime)
		defer cancel()
		sleep(ctx, sleepyTime)
		took := time.Since(startTime)
		assert.True(t, took < sleepyTime)
		assert.True(t, took > cancelTime)
	})

	t.Run("doesn't sleep with pre-canceled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		sleepyTime := 2 * time.Millisecond
		startTime := time.Now()
		sleep(ctx, sleepyTime)
		took := time.Since(startTime)
		assert.True(t, took < sleepyTime)
	})
}
