package retry

import "context"

type OptionFunc func(*Options)

type OnRetryFunc func(attempt uint, err error)

type Options struct {
	ctx context.Context

	// Default is 3
	maxRetries uint

	// Default is 3 seconds.
	backoffFunc BackoffFunc

	onRetryFunc OnRetryFunc
}

func WithContext(ctx context.Context) OptionFunc {
	return func(opt *Options) {
		opt.ctx = ctx
	}
}

func MaxRetries(maxRetries uint) OptionFunc {
	return func(opt *Options) {
		opt.maxRetries = maxRetries
	}
}

func WithBackoff(backoffFunc BackoffFunc) OptionFunc {
	return func(opt *Options) {
		opt.backoffFunc = backoffFunc
	}
}

func OnRetry(onRetryFunc OnRetryFunc) OptionFunc {
	return func(opt *Options) {
		opt.onRetryFunc = onRetryFunc
	}
}
