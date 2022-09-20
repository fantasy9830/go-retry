package goretry

type OptionFunc func(*options)

type options struct {
	// Default is 3
	maxRetries uint

	// Default is 3 seconds.
	backoffFunc BackoffFunc
}

func MaxRetries(maxRetries uint) OptionFunc {
	return func(opt *options) {
		opt.maxRetries = maxRetries
	}
}

func WithBackoff(backoffFunc BackoffFunc) OptionFunc {
	return func(opt *options) {
		opt.backoffFunc = backoffFunc
	}
}
