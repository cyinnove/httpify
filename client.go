package httpify

import (
	"net/http"
	"time"
)

// Client wraps an HTTP client with retry and retryStrategy logic.
type Client struct {
	HTTPClient      *http.Client
	RequestLogHook  RequestLogHook
	ResponseLogHook ResponseLogHook
	ErrorHandler    ErrorHandler
	CheckRetry      CheckRetry
	RetryStrategy   RetryStrategy
	options         Options
}

// Options defines retryable settings for the HTTP client.
type Options struct {
	RetryWaitMin  time.Duration
	RetryWaitMax  time.Duration
	Timeout       time.Duration
	RetryMax      int
	RespReadLimit int64
	KillIdleConn  bool
}

// Default options for spraying multiple hosts.
var DefaultOptionsSpraying = Options{
	RetryWaitMin:  1 * time.Second,
	RetryWaitMax:  30 * time.Second,
	Timeout:       30 * time.Second,
	RetryMax:      5,
	RespReadLimit: 4096,
	KillIdleConn:  true,
}

// Default options for targeting a single host.
var DefaultOptionsSingle = Options{
	RetryWaitMin:  1 * time.Second,
	RetryWaitMax:  30 * time.Second,
	Timeout:       30 * time.Second,
	RetryMax:      5,
	RespReadLimit: 4096,
	KillIdleConn:  false,
}

// NewClient initializes a Client with specified options.
func NewClient(options Options) *Client {
	httpClient := DefaultHTTPClient(options.Timeout)
	return &Client{
		HTTPClient:    httpClient,
		CheckRetry:    DefaultRetryPolicy(),
		RetryStrategy: DefaultRetryStrategy(),
		options:       options,
	}
}

// NewWithHTTPClient initializes a Client with a custom HTTP client.
func NewWithHTTPClient(client *http.Client, options Options) *Client {
	return &Client{
		HTTPClient:    client,
		CheckRetry:    DefaultRetryPolicy(),
		RetryStrategy: DefaultRetryStrategy(),
		options:       options,
	}
}

// DefaultHTTPClient creates an HTTP client with a default timeout.
func DefaultHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{Timeout: timeout, Transport: NoKeepAliveTransport()}
}

// setKillIdleConnections configures connection keep-alive behavior based on options.
func (c *Client) setKillIdleConnections() {
	if c.HTTPClient != nil || !c.options.KillIdleConn {
		if b, ok := c.HTTPClient.Transport.(*http.Transport); ok {
			c.options.KillIdleConn = b.DisableKeepAlives || b.MaxConnsPerHost < 0
		}
	}
}
