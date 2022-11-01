package goretry_test

import (
	"context"
	"errors"

	"testing"
	"time"

	goretry "github.com/fantasy9830/go-retry"
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

func TestGoRetry(t *testing.T) {
	suite.Run(t, new(Suit))
}

func (s *Suit) TestDefault() {
	var retryCount uint
	maxRetries := uint(3)

	start := time.Now()
	err := goretry.Do(func(ctx context.Context) error {
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

	options := []goretry.OptionFunc{
		goretry.WithBackoff(goretry.BackoffLinear(s.Duration)),
	}

	start := time.Now()
	err := goretry.Do(func(ctx context.Context) error {
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
		var retryCount int
		maxRetries := 10

		options := []goretry.OptionFunc{
			goretry.MaxRetries(0),
			goretry.WithBackoff(goretry.BackoffLinear(s.Duration)),
		}

		err := goretry.Do(func(ctx context.Context) error {
			retryCount++
			if retryCount == maxRetries {
				return nil
			}

			return errors.New("TestMaxRetries")
		}, options...)

		assert.Nil(t, err)
		assert.EqualValues(s.T(), maxRetries, retryCount)
	})

	s.T().Run("MaxRetries is 10", func(t *testing.T) {
		var retryCount uint
		maxRetries := uint(10)

		options := []goretry.OptionFunc{
			goretry.MaxRetries(maxRetries),
			goretry.WithBackoff(goretry.BackoffLinear(s.Duration)),
		}

		err := goretry.Do(func(ctx context.Context) error {
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
		options := []goretry.OptionFunc{
			goretry.WithBackoff(goretry.BackoffLinear(s.Duration)),
		}

		start := time.Now()
		err := goretry.Do(func(ctx context.Context) error {
			return errors.New("BackoffLinear")
		}, options...)
		duration := time.Since(start)

		assert.Error(t, err)
		assert.Equal(t, "BackoffLinear", err.Error())
		assert.WithinRange(t, start.Add(duration), start, start.Add(3*s.Duration))
	})

	s.T().Run("BackoffExponential", func(t *testing.T) {
		options := []goretry.OptionFunc{
			goretry.MaxRetries(4),
			goretry.WithBackoff(goretry.BackoffExponential(s.Duration)),
		}

		start := time.Now()
		err := goretry.Do(func(ctx context.Context) error {
			return errors.New("BackoffExponential")
		}, options...)
		duration := time.Since(start)

		assert.Error(t, err)
		assert.Equal(t, "BackoffExponential", err.Error())
		assert.WithinRange(t, start.Add(duration), start, start.Add(8*s.Duration))
	})
}

func (s *Suit) TestContext() {
	s.T().Run("context with cancel", func(t *testing.T) {
		c, cancel := context.WithCancel(context.Background())
		retrySum := 0

		options := []goretry.OptionFunc{
			goretry.WithContext(c),
			goretry.MaxRetries(0),
			goretry.WithBackoff(goretry.BackoffExponential(s.Duration)),
		}

		err := goretry.Do(func(ctx context.Context) error {
			retrySum++
			if retrySum == 2 {
				cancel()
			}

			return errors.New("TestContext")
		}, options...)

		assert.Error(t, err)
		assert.Equal(t, "context canceled", err.Error())
		assert.EqualValues(t, 2, retrySum)
	})

	s.T().Run("context with deadline", func(t *testing.T) {
		c, cancel := context.WithTimeout(context.Background(), 2*s.Duration)
		defer cancel()

		options := []goretry.OptionFunc{
			goretry.WithContext(c),
		}

		start := time.Now()
		err := goretry.Do(func(ctx context.Context) error {
			return errors.New("TestContext")
		}, options...)
		duration := time.Since(start)

		assert.Error(t, err)
		assert.Equal(t, "context deadline exceeded", err.Error())
		assert.WithinRange(t, start.Add(duration), start, start.Add(3*s.Duration))
	})
}
