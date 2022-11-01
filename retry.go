package goretry

import (
	"context"
	"time"
)

type RetryableFunc func(ctx context.Context) error

func Do(retryableFunc RetryableFunc, optFuncs ...OptionFunc) (lastErr error) {
	// default options
	opt := &Options{
		ctx:         context.Background(),
		maxRetries:  3,
		backoffFunc: BackoffLinear(3 * time.Second),
	}

	for _, applyFunc := range optFuncs {
		applyFunc(opt)
	}

	if opt.maxRetries == 0 {
		for attempt := uint(0); ; attempt++ {
			if err := waitRetryBackoff(attempt, opt); err != nil {
				return err
			}

			lastErr = retryableFunc(opt.ctx)
			if lastErr == nil {
				return nil
			}
		}
	} else {
		for attempt := uint(0); attempt < opt.maxRetries; attempt++ {
			if err := waitRetryBackoff(attempt, opt); err != nil {
				return err
			}

			lastErr = retryableFunc(opt.ctx)
			if lastErr == nil {
				return nil
			}
		}
	}

	return lastErr
}

func waitRetryBackoff(attempt uint, opt *Options) (err error) {
	var waitTime time.Duration = 0
	if attempt > 0 {
		waitTime = opt.backoffFunc(attempt)
	}

	if waitTime > 0 {
		err = sleep(opt.ctx, waitTime)
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
