package httpify

import (
	"net/http"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	options := Options{
		RetryWaitMin: 1 * time.Second,
		RetryWaitMax: 5 * time.Second,
		Timeout:      10 * time.Second,
		RetryMax:     3,
	}
	client := NewClient(options)

	assert.NotNil(t, client)
	assert.Equal(t, client.options, options)
	assert.NotNil(t, client.HTTPClient)
}

func TestNewWithHTTPClient(t *testing.T) {
	options := Options{
		RetryWaitMin: 1 * time.Second,
		RetryWaitMax: 5 * time.Second,
		Timeout:      10 * time.Second,
		RetryMax:     3,
	}
	customClient := &http.Client{Timeout: 20 * time.Second}
	client := NewWithHTTPClient(customClient, options)

	assert.NotNil(t, client)
	assert.Equal(t, client.HTTPClient, customClient)
	assert.Equal(t, client.options, options)
}

func TestSetKillIdleConnections(t *testing.T) {
	options := Options{
		KillIdleConn: true,
	}

	// Initialize the client with a default transport
	client := NewClient(options)

	// Check if the transport is of type *http.Transport
	transport, ok := client.HTTPClient.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("Transport is not of expected type *http.Transport, got %T", client.HTTPClient.Transport)
	}

	// Set transport properties
	client.setKillIdleConnections()

	// Validate that idle connections are disabled as per options
	assert.True(t, transport.DisableKeepAlives, "Expected DisableKeepAlives to be true")
}
