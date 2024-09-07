package httpify

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultBackoff(t *testing.T) {
	backoff := DefaultBackoff()
	min := 1 * time.Second
	max := 10 * time.Second

	duration := backoff(min, max, 1, nil)
	assert.LessOrEqual(t, duration, max)
	assert.GreaterOrEqual(t, duration, min)
}

func TestLinearJitterBackoff(t *testing.T) {
	backoff := LinearJitterBackoff()
	min := 1 * time.Second
	max := 5 * time.Second

	for attemptNum := 1; attemptNum <= 10; attemptNum++ {
		duration := backoff(min, max, attemptNum, nil)
		assert.LessOrEqual(t, duration, max, "Backoff duration exceeded max value")
		assert.GreaterOrEqual(t, duration, min, "Backoff duration is less than min value")
	}
}


func TestFullJitterBackoff(t *testing.T) {
	backoff := FullJitterBackoff()
	min := 1 * time.Second
	max := 10 * time.Second

	duration := backoff(min, max, 1, nil)
	assert.LessOrEqual(t, duration, max)
	assert.GreaterOrEqual(t, duration, min)
}

func TestExponentialJitterBackoff(t *testing.T) {
	backoff := ExponentialJitterBackoff()
	min := 1 * time.Second
	max := 10 * time.Second

	duration := backoff(min, max, 1, nil)
	assert.LessOrEqual(t, duration, max)
	assert.GreaterOrEqual(t, duration, min)
}
