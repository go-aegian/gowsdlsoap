package soap

import "fmt"

// HTTPError is returned whenever the HTTP request to the server fails
type HTTPError struct {
	StatusCode   int
	ResponseBody []byte
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP Status %d: %s", e.StatusCode, string(e.ResponseBody))
}
