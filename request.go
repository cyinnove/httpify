package httpify

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"os"
)

// LenReader interface defines a method to get the length of a reader.
type LenReader interface {
	Len() int
}

// Request wraps HTTP request metadata for retries.
type Request struct {
	body ReaderFunc
	*http.Request
	Metrics Metrics
}

// Metrics stores retry and error metrics for a request.
type Metrics struct {
	Failures    int
	Retries     int
	DrainErrors int
}

// RequestLogHook allows executing custom logic before each retry.
type RequestLogHook func(*http.Request, int)

// ResponseLogHook allows executing custom logic after each HTTP request.
type ResponseLogHook func(*http.Response)

// ErrorHandler handles retries exhaustion and returns custom responses.
type ErrorHandler func(resp *http.Response, err error, numTries int) (*http.Response, error)

// ReaderFunc defines a function type for creating readers.
type ReaderFunc func() (io.Reader, error)

// NewRequest creates a new wrapped request.
func NewRequest(method, url string, body interface{}) (*Request, error) {
	bodyReader, contentLength, err := getBodyReaderAndContentLength(body)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.ContentLength = contentLength

	return &Request{bodyReader, httpReq, Metrics{}}, nil
}

// NewRequestWithContext creates a new wrapped request with a context.
func NewRequestWithContext(ctx context.Context, method, url string, body interface{}) (*Request, error) {
	bodyReader, contentLength, err := getBodyReaderAndContentLength(body)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	httpReq.ContentLength = contentLength

	return &Request{bodyReader, httpReq, Metrics{}}, nil
}

// WithContext returns a shallow copy of the request with a new context.
func (r *Request) WithContext(ctx context.Context) *Request {
	r.Request = r.Request.WithContext(ctx)
	return r
}

// FromRequest wraps an http.Request into a retryable Request.
func FromRequest(r *http.Request) (*Request, error) {
	bodyReader, contentLength, err := getBodyReaderAndContentLength(r.Body)
	if err != nil {
		return nil, err
	}
	r.ContentLength = contentLength

	// Reset the body on the original request
	r.Body = io.NopCloser(bytes.NewReader([]byte{}))

	return &Request{
		body:    bodyReader,
		Request: r,
		Metrics: Metrics{},
	}, nil
}


// FromRequestWithTrace wraps an http.Request into a retryable Request with trace enabled.
func FromRequestWithTrace(r *http.Request) (*Request, error) {
	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Fprintf(os.Stderr, "Got connection (Reused: %v, Idle: %v, IdleTime: %v)\n", connInfo.Reused, connInfo.WasIdle, connInfo.IdleTime)
		},
		ConnectStart: func(network, addr string) {
			fmt.Fprintf(os.Stderr, "Connecting (network: %s, address: %s)\n", network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			fmt.Fprintf(os.Stderr, "Connected (network: %s, address: %s, error: %v)\n", network, addr, err)
		},
		GotFirstResponseByte: func() {
			fmt.Fprintf(os.Stderr, "Received first byte of response\n")
		},
		WroteHeaders: func() {
			fmt.Fprintf(os.Stderr, "Request headers written\n")
		},
		WroteRequest: func(wr httptrace.WroteRequestInfo) {
			fmt.Fprintf(os.Stderr, "Request sent (error: %v)\n", wr.Err)
		},
	}

	r = r.WithContext(httptrace.WithClientTrace(r.Context(), trace))
	return FromRequest(r)
}

// BodyBytes returns a copy of the request body data.
func (r *Request) BodyBytes() ([]byte, error) {
	if r.body == nil {
		return nil, nil
	}
	body, err := r.body()
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(body)
	return buf.Bytes(), err
}
func getBodyReaderAndContentLength(rawBody interface{}) (ReaderFunc, int64, error) {
	var bodyReader ReaderFunc
	var contentLength int64

	if rawBody != nil {
		switch body := rawBody.(type) {
		case *bytes.Buffer:
			buf := body.Bytes()
			bodyReader = func() (io.Reader, error) { return bytes.NewReader(buf), nil }
			contentLength = int64(len(buf))
		case *bytes.Reader:
			buf, err := io.ReadAll(body)
			if err != nil {
				return nil, 0, err
			}
			bodyReader = func() (io.Reader, error) { return bytes.NewReader(buf), nil }
			contentLength = int64(len(buf))
		case io.Reader:
			buf, err := io.ReadAll(body)
			if err != nil {
				return nil, 0, err
			}
			bodyReader = func() (io.Reader, error) { return bytes.NewReader(buf), nil }
			contentLength = int64(len(buf))
		default:
			return nil, 0, nil
		}
	}

	return bodyReader, contentLength, nil
}


func createReaderAndGetLength(body ReaderFunc) (int64, bool) {
	tmp, err := body()
	if err != nil {
		return 0, false
	}
	if lr, ok := tmp.(LenReader); ok {
		return int64(lr.Len()), true
	}
	if c, ok := tmp.(io.Closer); ok {
		c.Close()
	}
	return 0, false
}
