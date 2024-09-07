package httpify

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultRetryStrategy(t *testing.T) {
	retryStrategy := DefaultRetryStrategy()
	min := 1 * time.Second
	max := 10 * time.Second

	duration := retryStrategy(min, max, 1, nil)
	assert.LessOrEqual(t, duration, max)
	assert.GreaterOrEqual(t, duration, min)
}

func TestLinearRandomizedRetryStrategy(t *testing.T) {
	retryStrategy := LinearRandomizedRetryStrategy()
	min := 1 * time.Second
	max := 5 * time.Second

	for attemptNum := 1; attemptNum <= 10; attemptNum++ {
		duration := retryStrategy(min, max, attemptNum, nil)
		assert.LessOrEqual(t, duration, max, "RetryStrategy duration exceeded max value")
		assert.GreaterOrEqual(t, duration, min, "RetryStrategy duration is less than min value")
	}
}

func TestRandomizedFullRetryStrategy(t *testing.T) {
	retryStrategy := RandomizedFullRetryStrategy()
	min := 1 * time.Second
	max := 10 * time.Second

	duration := retryStrategy(min, max, 1, nil)
	assert.LessOrEqual(t, duration, max)
	assert.GreaterOrEqual(t, duration, min)
}

func TestExponentialRandomizedRetryStrategy(t *testing.T) {
	retryStrategy := ExponentialRandomizedRetryStrategy()
	min := 1 * time.Second
	max := 10 * time.Second

	duration := retryStrategy(min, max, 1, nil)
	assert.LessOrEqual(t, duration, max)
	assert.GreaterOrEqual(t, duration, min)
}
