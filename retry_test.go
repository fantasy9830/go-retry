package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fantasy9830/go-retry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Suit struct {
	suite.Suite
	Duration time.Duration
}

func (s *Suit) SetupSuite() {
	s.Duration = 100 * time.Millisecond
}

func TestRetry(t *testing.T) {
	suite.Run(t, new(Suit))
}

func (s *Suit) TestDefault() {
	var retryCount uint
	maxRetries := uint(3)

	start := time.Now()
	err := retry.Do(func(ctx context.Context) error {
		retryCount++
		return errors.New("TestDefault")
	})
	duration := time.Since(start)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), "TestDefault", err.Error())
	assert.WithinRange(s.T(), start.Add(duration), start.Add(6*time.Second), start.Add(7*time.Second))
	assert.EqualValues(s.T(), maxRetries, retryCount)
}

func (s *Suit) TestNoError() {
	var retryCount uint
	maxRetries := uint(1)

	options := []retry.OptionFunc{
		retry.WithBackoff(retry.BackoffLinear(s.Duration)),
	}

	start := time.Now()
	err := retry.Do(func(ctx context.Context) error {
		retryCount++
		return nil
	}, options...)
	duration := time.Since(start)

	assert.Nil(s.T(), err)
	assert.WithinRange(s.T(), start.Add(duration), start, start.Add(2*s.Duration))
	assert.EqualValues(s.T(), maxRetries, retryCount)
}

func (s *Suit) TestMaxRetries() {
	s.T().Run("MaxRetries is 0", func(t *testing.T) {
		var retryCount uint
		var retryError, nilError error
		maxRetries := uint(10)

		options := []retry.OptionFunc{
			retry.MaxRetries(0),
			retry.WithBackoff(retry.BackoffLinear(s.Duration)),
			retry.OnRetry(func(attempt uint, err error) {
				retryCount = attempt
				if err != nil {
					retryError = err
				}
				nilError = err
			}),
		}

		err := retry.Do(func(ctx context.Context) error {
			if retryCount == maxRetries-1 {
				return nil
			}

			return errors.New("TestMaxRetries")
		}, options...)

		assert.Nil(t, err)
		assert.EqualValues(t, maxRetries, retryCount)
		assert.Equal(t, "TestMaxRetries", retryError.Error())
		assert.Nil(t, nilError)
	})

	s.T().Run("MaxRetries is 10", func(t *testing.T) {
		var retryCount uint
		maxRetries := uint(10)

		options := []retry.OptionFunc{
			retry.MaxRetries(maxRetries),
			retry.WithBackoff(retry.BackoffLinear(s.Duration)),
		}

		err := retry.Do(func(ctx context.Context) error {
			retryCount++
			return errors.New("TestMaxRetries")
		}, options...)

		assert.Error(s.T(), err)
		assert.Equal(s.T(), "TestMaxRetries", err.Error())
		assert.EqualValues(s.T(), maxRetries, retryCount)
	})
}

func (s *Suit) TestWithBackoff() {
	s.T().Run("BackoffLinear", func(t *testing.T) {
		options := []retry.OptionFunc{
			retry.WithBackoff(retry.BackoffLinear(s.Duration)),
		}

		start := time.Now()
		err := retry.Do(func(ctx context.Context) error {
			return errors.New("BackoffLinear")
		}, options...)
		duration := time.Since(start)

		assert.Error(t, err)
		assert.Equal(t, "BackoffLinear", err.Error())
		assert.WithinRange(t, start.Add(duration), start, start.Add(3*s.Duration))
	})

	s.T().Run("BackoffExponential", func(t *testing.T) {
		options := []retry.OptionFunc{
			retry.MaxRetries(4),
			retry.WithBackoff(retry.BackoffExponential(s.Duration)),
		}

		start := time.Now()
		err := retry.Do(func(ctx context.Context) error {
			return errors.New("BackoffExponential")
		}, options...)
		duration := time.Since(start)

		assert.Error(t, err)
		assert.Equal(t, "BackoffExponential", err.Error())
		assert.WithinRange(t, start.Add(duration), start, start.Add(8*s.Duration))
	})
}

func (s *Suit) TestWithContext() {
	s.T().Run("context with cancel", func(t *testing.T) {
		c, cancel := context.WithCancel(context.Background())
		retrySum := 0

		options := []retry.OptionFunc{
			retry.WithContext(c),
			retry.MaxRetries(0),
			retry.WithBackoff(retry.BackoffExponential(s.Duration)),
		}

		err := retry.Do(func(ctx context.Context) error {
			retrySum++
			if retrySum == 2 {
				cancel()
			}

			return errors.New("TestContext")
		}, options...)

		assert.Error(t, err)
		assert.Equal(t, "TestContext", err.Error())
		assert.EqualValues(t, 2, retrySum)
	})

	s.T().Run("context with deadline", func(t *testing.T) {
		c, cancel := context.WithTimeout(context.Background(), 2*s.Duration)
		defer cancel()

		options := []retry.OptionFunc{
			retry.WithContext(c),
		}

		start := time.Now()
		err := retry.Do(func(ctx context.Context) error {
			return errors.New("TestContext")
		}, options...)
		duration := time.Since(start)

		assert.Error(t, err)
		assert.Equal(t, "TestContext", err.Error())
		assert.WithinRange(t, start.Add(duration), start, start.Add(3*s.Duration))
	})
}

func (s *Suit) TestOnRetry() {
	s.T().Run("OnRetry", func(t *testing.T) {
		var retryCount uint
		var retryError error
		maxRetries := uint(3)
		options := []retry.OptionFunc{
			retry.MaxRetries(maxRetries),
			retry.WithBackoff(retry.BackoffLinear(s.Duration)),
			retry.OnRetry(func(attempt uint, err error) {
				retryCount = attempt
				retryError = err
			}),
		}

		err := retry.Do(func(ctx context.Context) error {
			return errors.New("OnRetry")
		}, options...)

		assert.Error(t, err)
		assert.Equal(t, "OnRetry", err.Error())
		assert.Equal(t, "OnRetry", retryError.Error())
		assert.EqualValues(s.T(), maxRetries, retryCount)
	})
}
