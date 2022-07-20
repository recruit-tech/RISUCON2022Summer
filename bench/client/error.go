package client

import (
	"errors"
	"fmt"
)

var (
	ErrNoContent            = errors.New("204 No Content")
	ErrBadRequest           = errors.New("400 Bad Request")
	ErrUnauthorized         = errors.New("401 Unauthorized")
	ErrNotFound             = errors.New("404 Not Found")
	ErrConflict             = errors.New("409 Conflict")
	ErrServiceUnavailable   = errors.New("503 Service Unavailable")
	ErrUnexpectedStatusCode = errors.New("予期せぬstatus codeです")
)

func newUnexpectedStatusCodeError(endpoint string, statusCode int) error {
	return fmt.Errorf("%s: %w (%d)", endpoint, ErrUnexpectedStatusCode, statusCode)
}
