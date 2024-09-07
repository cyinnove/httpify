package httpify

import (
	"io"
	"net/http"
)

// Discard discards the response body and closes the underlying connection.
// It reads up to RespReadLimit bytes from the response body to ensure connection reuse.
// If an error occurs during reading, it increments the DrainErrors metric.
func Discard(req *Request, resp *http.Response, respReadLimit int64) {
	defer resp.Body.Close()

	_, err := io.Copy(io.Discard, io.LimitReader(resp.Body, respReadLimit))
	if err != nil {
		req.Metrics.DrainErrors++
	}
}
