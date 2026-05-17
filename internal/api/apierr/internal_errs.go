// Package apierr indented for viewable and readable API errors
package apierr

import "net/http"

var (
	codeBadRequest      = "bad_request"
	codeValidationError = "validation_error"
)

// Error it's a struct indented for viewable API error
type Error struct {
	RequestID string `json:"request_id"`
	Status    int    `json:"status_code"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

// BadRequestError it's a constructor for generic http.StatusBadRequest API errors (ex.: incorrect JSON body)
func BadRequestError(message, requestID string) Error {
	return Error{
		Status:    http.StatusBadRequest,
		RequestID: requestID,
		Code:      codeBadRequest,
		Message:   message,
	}
}

// ValidationFailedError it's a constructor for validation http.StatusBadRequest API errors. See `validate` tags in package `entities/dto`
func ValidationFailedError(message, requestID string) Error {
	return Error{
		Status:    http.StatusBadRequest,
		RequestID: requestID,
		Code:      codeValidationError,
		Message:   message,
	}
}
