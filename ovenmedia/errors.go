package ovenmedia

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError is returned when OME replies with a non-2xx status. OME's error
// replies carry {"statusCode", "message"}; Body holds the raw reply.
type APIError struct {
	StatusCode int
	Message    string
	Body       []byte
}

func (e *APIError) Error() string {
	return fmt.Sprintf("ovenmedia: HTTP %d: %s", e.StatusCode, e.Message)
}

func newAPIError(status int, body []byte) *APIError {
	e := &APIError{StatusCode: status, Body: body}
	var r BaseResponseOK
	if json.Unmarshal(body, &r) == nil && r.Message != "" {
		e.Message = r.Message
	} else {
		e.Message = http.StatusText(status)
	}
	return e
}
