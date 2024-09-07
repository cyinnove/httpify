package httpify

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultHostSprayingTransport(t *testing.T) {
	transport := NoKeepAliveTransport()
	assert.NotNil(t, transport)
	assert.True(t, transport.DisableKeepAlives)
}

func TestDefaultReusePooledTransport(t *testing.T) {
	transport := PooledTransport()
	assert.NotNil(t, transport)
	assert.Equal(t, 100, transport.MaxIdleConnsPerHost)
	assert.Equal(t, 90*time.Second, transport.IdleConnTimeout)
}

func TestDefaultClient(t *testing.T) {
	client := DefaultClient()
	assert.NotNil(t, client)
	assert.NotNil(t, client.Transport)
}

func TestDefaultPooledClient(t *testing.T) {
	client := DefaultPooledClient()
	assert.NotNil(t, client)
	assert.NotNil(t, client.Transport)
}
