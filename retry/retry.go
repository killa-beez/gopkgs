package retry

import (
	"context"
	"time"
)

const defaultMaxAttempts = 3

func defaultSleepInterval(int, error) time.Duration {
	return 0
}

func defaultRetryable(lastErr error) bool {
	return lastErr != nil
}

//Config configuration for Retry
type Config struct {
	MaxAttempts int                                               // How many attempts to make before giving up. Default is 3
	Backoff     func(retryCount int, lastErr error) time.Duration // function deciding long to wait before the next retry. Default always returns 0.
	Retryable   func(lastErr error) bool                          // function deciding whether an error is retryable. Default is all non-nil errors are retryable.
}

/*
Retry keeps running fn() until one of:
 - it doesn't return an error
 - Retryable(lastErr) returns false
 - ctx is canceled
 - config.MaxAttempts is exceeded

Retry always returns the last error returned by fn().
*/
func Retry(ctx context.Context, fn func() error, config Config) error {
	maxAttempts := config.MaxAttempts
	if maxAttempts == 0 {
		maxAttempts = defaultMaxAttempts
	}
	sleepInterval := config.Backoff
	if sleepInterval == nil {
		sleepInterval = defaultSleepInterval
	}
	retryable := config.Retryable
	if retryable == nil {
		retryable = defaultRetryable
	}
	var err error
	count := 0
	for {
		count++
		if count > maxAttempts {
			break
		}
		err = fn()
		if !retryable(err) {
			break
		}
		interval := sleepInterval(count, err)
		sleep(ctx, interval)
		if ctx.Err() != nil {
			break
		}
	}
	return err
}

//sleep returns when either dur expires or ctx is canceled
func sleep(ctx context.Context, dur time.Duration) {
	select {
	case <-time.After(dur):
	case <-ctx.Done():
	}
}
