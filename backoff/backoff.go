package backoff

import (
	"context"
	"math"
	"math/rand"
	"time"
)

type (
	//Backoff is like a simplified, stateless version of cenk/backoff
	Backoff interface {
		Interval(retry int) time.Duration
		Sleep(ctx context.Context, retry int)
		Retry(ctx context.Context, retries int, fn func() error) error
	}

	backoff struct {
		multiplier   float64
		baseDuration float64
		maxDuration  float64
		jitter       bool
	}
)

//NewBackoff returns a new Backoff
func NewBackoff(multiplier float64, baseDuration, maxDuration time.Duration, jitter bool) Backoff {
	return &backoff{
		multiplier:   multiplier,
		baseDuration: float64(baseDuration),
		maxDuration:  float64(maxDuration),
		jitter:       jitter,
	}
}

func (b *backoff) Interval(retry int) time.Duration {
	interval := b.baseDuration * math.Pow(b.multiplier, float64(retry))
	if b.jitter {
		interval += (b.baseDuration * rand.NormFloat64())
	}
	if b.maxDuration > 0.0 && interval > b.maxDuration {
		interval = b.maxDuration
	}
	if b.baseDuration > interval {
		interval = b.baseDuration
	}
	return time.Duration(interval)
}

func (b *backoff) Sleep(ctx context.Context, retry int) {
	select {
	case <-time.After(b.Interval(retry)):
	case <-ctx.Done():
	}
}

//Retry keeps retrying fn until one of:
// fn returns nil
// retry count is equal to retries
// ctx is canceled
//
// if retries is less than 0, there is no max number or retries
//Retry returns the last return value of fn whether it's error or nil
func (b *backoff) Retry(ctx context.Context, retries int, fn func() error) error {
	var err error
	for i := 0; retries != i; i++ {
		err = fn()
		if err == nil {
			break
		}
		b.Sleep(ctx, i+1)
		if ctx.Err() != nil {
			break
		}
	}
	return err
}
