package retry

import (
	"time"
)

type BackoffFunc func(attempt uint) time.Duration

func BackoffLinear(duration time.Duration) BackoffFunc {
	return func(attempt uint) time.Duration {
		return duration
	}
}

func BackoffExponential(duration time.Duration) BackoffFunc {
	return func(attempt uint) time.Duration {
		// 2^(attempt-1)
		return ((1 << attempt) >> 1) * duration
	}
}
