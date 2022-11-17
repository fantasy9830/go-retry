package retry_test

import (
	"context"
	"testing"

	"github.com/fantasy9830/go-retry"
)

func BenchmarkRetryDo(b *testing.B) {
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		_ = retry.Do(func(ctx context.Context) error {
			return nil
		})
	}
}
