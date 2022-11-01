package goretry

import "context"

type OptionFunc func(*Options)

type Options struct {
	ctx context.Context

	// Default is 3
	maxRetries uint

	// Default is 3 seconds.
	backoffFunc BackoffFunc
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
