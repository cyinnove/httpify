package httpify

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

// DefaultHostSprayingTransport returns a new http.Transport with disabled idle connections and keepalives.
func DefaultHostSprayingTransport() *http.Transport {
	transport := DefaultReusePooledTransport()
	transport.DisableKeepAlives = true
	transport.MaxIdleConnsPerHost = -1
	return transport
}

// DefaultReusePooledTransport returns a new http.Transport for connection reuse.
func DefaultReusePooledTransport() *http.Transport {
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
		Transport: DefaultHostSprayingTransport(),
	}
}

// DefaultPooledClient returns an http.Client with a pooled transport for connection reuse.
func DefaultPooledClient() *http.Client {
	return &http.Client{
		Transport: DefaultReusePooledTransport(),
	}
}
