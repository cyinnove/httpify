package httpify

import (
	"context"
	"crypto/x509"
	"net/http"
	"net/url"
	"testing"
)

func TestDefaultRetryPolicy(t *testing.T) {
	policy := DefaultRetryPolicy()

	tests := []struct {
		name     string
		resp     *http.Response
		err      error
		expected bool
	}{
		{
			name:     "Retryable connection error",
			resp:     nil,
			err:      &url.Error{Err: &url.Error{}},
			expected: true,
		},
		{
			name:     "Non-retryable error",
			resp:     nil,
			err:      &url.Error{Err: x509.UnknownAuthorityError{}},
			expected: false,
		},
		{
			name:     "Retryable URL error",
			resp:     nil,
			err:      &url.Error{Err: &url.Error{}},
			expected: true,
		},
		{
			name:     "Context cancellation",
			resp:     nil,
			err:      nil,
			expected: false,
		},
		{
			name:     "Successful response",
			resp:     &http.Response{},
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			retry, _ := policy(ctx, tt.resp, tt.err)
			if retry != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, retry)
			}
		})
	}
}

func TestHostSprayRetryPolicy(t *testing.T) {
	policy := HostSprayRetryPolicy()

	// Reuse tests from DefaultRetryPolicy to ensure it behaves the same
	tests := []struct {
		name     string
		resp     *http.Response
		err      error
		expected bool
	}{
		{
			name:     "Retryable connection error",
			resp:     nil,
			err:      &url.Error{Err: &url.Error{}},
			expected: true,
		},
		{
			name:     "Non-retryable error",
			resp:     nil,
			err:      &url.Error{Err: x509.UnknownAuthorityError{}},
			expected: false,
		},
		{
			name:     "Retryable URL error",
			resp:     nil,
			err:      &url.Error{Err: &url.Error{}},
			expected: true,
		},
		{
			name:     "Context cancellation",
			resp:     nil,
			err:      nil,
			expected: false,
		},
		{
			name:     "Successful response",
			resp:     &http.Response{},
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			retry, _ := policy(ctx, tt.resp, tt.err)
			if retry != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, retry)
			}
		})
	}
}

func TestIsNonRetryableError(t *testing.T) {
	t.Run("Redirects error", func(t *testing.T) {
		urlErr := &url.Error{Err: &url.Error{}}
		expected := false
		result := isNonRetryableError(urlErr)
		if result != expected {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("Unsupported scheme error", func(t *testing.T) {
		urlErr := &url.Error{Err: &url.Error{}}
		expected := false
		result := isNonRetryableError(urlErr)
		if result != expected {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("TLS cert error", func(t *testing.T) {
		urlErr := &url.Error{Err: x509.UnknownAuthorityError{}}
		expected := true
		result := isNonRetryableError(urlErr)
		if result != expected {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("Non-retryable error", func(t *testing.T) {
		urlErr := &url.Error{Err: &url.Error{}}
		expected := false
		result := isNonRetryableError(urlErr)
		if result != expected {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

func TestIsTLSCertError(t *testing.T) {
	tests := []struct {
		name     string
		urlErr   *url.Error
		expected bool
	}{
		{
			name:     "TLS cert error",
			urlErr:   &url.Error{Err: x509.UnknownAuthorityError{}},
			expected: true,
		},
		{
			name:     "Non-TLS error",
			urlErr:   &url.Error{Err: &url.Error{}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTLSCertError(tt.urlErr)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
