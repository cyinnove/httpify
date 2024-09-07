package httpify

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDoSuccess(t *testing.T) {
	options := Options{
		RetryMax: 3,
		Timeout:  10 * time.Second,
	}
	client := NewClient(options)
	req, _ := NewRequest(http.MethodGet, "http://example.com", nil)

	// Mocking HTTP Client for success case
	client.HTTPClient = &http.Client{
		Transport: &http.Transport{},
	}

	resp, err := client.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestDoWithRetry(t *testing.T) {
	options := Options{
		RetryMax: 3,
		Timeout:  5 * time.Second,
	}
	client := NewClient(options)
	req, _ := NewRequest(http.MethodGet, "http://invalid-url", nil)

	// Test retry logic by calling the "Do" method with a failing URL
	resp, err := client.Do(req)
	assert.NotNil(t, err)
	assert.Nil(t, resp)
}

