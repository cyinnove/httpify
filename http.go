package httpify

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

// HostSprayingTransport returns a new http.Transport with disabled idle connections and keepalives.
func NoKeepAliveTransport() *http.Transport {
	transport := PooledTransport()
	transport.DisableKeepAlives = true
	transport.MaxIdleConnsPerHost = -1
	return transport
}

// PooledTransport returns a new http.Transport for connection reuse.
func PooledTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:           100,
		IdleConnTimeout:        90 * time.Second,
		TLSHandshakeTimeout:    10 * time.Second,
		ExpectContinueTimeout:  1 * time.Second,
		MaxIdleConnsPerHost:    100,
		MaxResponseHeaderBytes: 4096, // Default is 10MB
		TLSClientConfig: &tls.Config{
			Renegotiation:      tls.RenegotiateOnceAsClient,
			InsecureSkipVerify: true, // Optional, but unsafe
		},
	}
}

// DefaultClient creates a new http.Client with disabled idle connections and keepalives.
func DefaultClient() *http.Client {
	return &http.Client{
		Transport: NoKeepAliveTransport(),
	}
}

// DefaultPooledClient returns an http.Client with a pooled transport for connection reuse.
func DefaultPooledClient() *http.Client {
	return &http.Client{
		Transport: PooledTransport(),
	}
}
