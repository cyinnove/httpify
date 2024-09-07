package httpify

import (
	"math"
	"math/rand"
	"net/http"
	"time"
)

// Backoff defines how long to wait between retries.
type Backoff func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration

// DefaultBackoff implements exponential backoff based on attempt count, bounded by min and max durations.
func DefaultBackoff() Backoff {
	return func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		sleep := time.Duration(math.Pow(2, float64(attemptNum)) * float64(min))
		if sleep > max {
			return max
		}
		return sleep
	}
}

// jitterBackoff returns a jittered duration within the range of min to max.
func jitterBackoff(min, max time.Duration, attemptNum int, randSource *rand.Rand) time.Duration {
	jitter := randSource.Float64() * float64(max-min)
	jitterDuration := time.Duration(int64(jitter) + int64(min))
	return jitterDuration * time.Duration(attemptNum+1)
}

// LinearJitterBackoff adds random jitter to linear backoff.
func LinearJitterBackoff() Backoff {
	randSource := newRandSource()
	return func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if max <= min {
			return min * time.Duration(attemptNum+1)
		}

		// Calculate jittered backoff
		jitteredDuration := jitterBackoff(min, max, attemptNum, randSource)

		// Ensure the jittered backoff does not exceed the max duration
		if jitteredDuration > max {
			return max
		}
		return jitteredDuration
	}
}


// FullJitterBackoff implements capped exponential backoff with full jitter.
func FullJitterBackoff() Backoff {
	randSource := newRandSource()
	return func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		duration := attemptNum * 1000000000 << 1
		jitter := randSource.Intn(duration-attemptNum) + int(min)
		if jitter > int(max) {
			return max
		}
		return time.Duration(jitter)
	}
}

// ExponentialJitterBackoff adds jitter to exponential backoff.
func ExponentialJitterBackoff() Backoff {
	randSource := newRandSource()
	return func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		base := math.Pow(2, float64(attemptNum)) * float64(min)
		jitter := randSource.Float64() * (base - float64(min))
		sleep := time.Duration(base + jitter)
		if sleep > max {
			return max
		}
		return sleep
	}
}

// newRandSource initializes a rand source with a mutex for concurrency safety.
func newRandSource() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
