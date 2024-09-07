package httpify

import (
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	client := NewClient(DefaultOptionsSingle)
	resp, err := client.Get("http://example.com")
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestHead(t *testing.T) {
	client := NewClient(DefaultOptionsSingle)
	resp, err := client.Head("http://example.com")
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestPost(t *testing.T) {
	client := NewClient(DefaultOptionsSingle)
	resp, err := client.Post("http://example.com", "application/json", strings.NewReader(`{"key":"value"}`))
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}

func TestPostForm(t *testing.T) {
	client := NewClient(DefaultOptionsSingle)
	data := url.Values{}
	data.Set("key", "value")

	resp, err := client.PostForm("http://example.com", data)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
}
