package httpify

import (
	"math"
	"math/rand"
	"net/http"
	"time"
)

// RetryStrategy defines how long to wait between retries.
type RetryStrategy func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration

// DefaultRetryStrategy implements exponential retryStrategy based on attempt count, bounded by min and max durations.
func DefaultRetryStrategy() RetryStrategy {
	return func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		sleep := time.Duration(math.Pow(2, float64(attemptNum)) * float64(min))
		if sleep > max {
			return max
		}
		return sleep
	}
}

// jitterRetryStrategy returns a Randomized duration within the range of min to max.
func jitterRetryStrategy(min, max time.Duration, attemptNum int, randSource *rand.Rand) time.Duration {
	randomized := randSource.Float64() * float64(max-min)
	jitterDuration := time.Duration(int64(randomized) + int64(min))
	return jitterDuration * time.Duration(attemptNum+1)
}

// LinearRandomizedRetryStrategy adds random randomized to linear retryStrategy.
func LinearRandomizedRetryStrategy() RetryStrategy {
	randSource := newRandSource()
	return func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		if max <= min {
			return min * time.Duration(attemptNum+1)
		}

		// Calculate Randomized retryStrategy
		RandomizedDuration := jitterRetryStrategy(min, max, attemptNum, randSource)

		// Ensure the Randomized retryStrategy does not exceed the max duration
		if RandomizedDuration > max {
			return max
		}
		return RandomizedDuration
	}
}

// RandomizedFullRetryStrategy implements capped exponential retryStrategy with full randomized.
func RandomizedFullRetryStrategy() RetryStrategy {
	randSource := newRandSource()
	return func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		duration := attemptNum * 1000000000 << 1
		randomized := randSource.Intn(duration-attemptNum) + int(min)
		if randomized > int(max) {
			return max
		}
		return time.Duration(randomized)
	}
}

// ExponentialRandomizedRetryStrategy adds randomized to exponential retryStrategy.
func ExponentialRandomizedRetryStrategy() RetryStrategy {
	randSource := newRandSource()
	return func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		base := math.Pow(2, float64(attemptNum)) * float64(min)
		randomized := randSource.Float64() * (base - float64(min))
		sleep := time.Duration(base + randomized)
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
