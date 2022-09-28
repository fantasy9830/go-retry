package goretry

import (
	"context"
	"time"
)

type RetryableFunc func() error

func Do(ctx context.Context, retryableFunc RetryableFunc, optFuncs ...OptionFunc) (lastErr error) {
	// default options
	opt := &options{
		maxRetries:  3,
		backoffFunc: BackoffLinear(3 * time.Second),
	}

	for _, applyFunc := range optFuncs {
		applyFunc(opt)
	}

	if opt.maxRetries == 0 {
		for attempt := uint(0); ; attempt++ {
			if err := waitRetryBackoff(ctx, attempt, opt); err != nil {
				return err
			}

			lastErr = retryableFunc()
			if lastErr == nil {
				return nil
			}
		}
	} else {
		for attempt := uint(0); attempt < opt.maxRetries; attempt++ {
			if err := waitRetryBackoff(ctx, attempt, opt); err != nil {
				return err
			}

			lastErr = retryableFunc()
			if lastErr == nil {
				return nil
			}
		}
	}

	return lastErr
}

func waitRetryBackoff(ctx context.Context, attempt uint, opt *options) (err error) {
	var waitTime time.Duration = 0
	if attempt > 0 {
		waitTime = opt.backoffFunc(attempt)
	}

	if waitTime > 0 {
		err = sleep(ctx, waitTime)
	}

	return err
}

func sleep(ctx context.Context, waitTime time.Duration) error {
	timer := time.NewTimer(waitTime)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
