package httpify

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRequest(t *testing.T) {
	req, err := NewRequest(http.MethodGet, "http://example.com", nil)
	assert.Nil(t, err)
	assert.NotNil(t, req)
}

func TestNewRequestWithContext(t *testing.T) {
	ctx := context.Background()
	req, err := NewRequestWithContext(ctx, http.MethodPost, "http://example.com", nil)
	assert.Nil(t, err)
	assert.NotNil(t, req)
}

func TestBodyBytes(t *testing.T) {
	bodyContent := []byte("test body")
	req, _ := NewRequest(http.MethodPost, "http://example.com", bodyContent)

	bodyBytes, err := req.BodyBytes()
	assert.Nil(t, err)
	assert.Equal(t, bodyContent, bodyBytes)
}
