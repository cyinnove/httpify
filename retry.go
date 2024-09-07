package httpify

import (
	"context"
	"crypto/x509"
	"net/http"
	"net/url"
	"regexp"
)

var (
	redirectsErrorRegex = regexp.MustCompile(`stopped after \d+ redirects\z`)
	schemeErrorRegex    = regexp.MustCompile(`unsupported protocol scheme`)
)

// CheckRetry defines a policy for retrying requests based on the response and error.
type CheckRetry func(ctx context.Context, resp *http.Response, err error) (bool, error)

// DefaultRetryPolicy retries on connection errors and server errors.
func DefaultRetryPolicy() CheckRetry {
	return func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if ctx.Err() != nil {
			return false, ctx.Err()
		}

		if err != nil {
			if urlErr, ok := err.(*url.Error); ok {
				// Handle specific error conditions
				if isNonRetryableError(urlErr) {
					return false, nil
				}
			}
			return true, nil // Retry on likely recoverable error
		}

		return false, nil
	}
}

// HostSprayRetryPolicy retries on connection and server errors for host-spraying use cases.
func HostSprayRetryPolicy() CheckRetry {
	return DefaultRetryPolicy()
}

func isNonRetryableError(urlErr *url.Error) bool {
	return redirectsErrorRegex.MatchString(urlErr.Error()) ||
		schemeErrorRegex.MatchString(urlErr.Error()) ||
		isTLSCertError(urlErr)
}

func isTLSCertError(urlErr *url.Error) bool {
	_, ok := urlErr.Err.(x509.UnknownAuthorityError)
	return ok
}
