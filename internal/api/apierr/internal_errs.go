package apierr

import "net/http"

var (
	CodeBadRequest      = "bad_request"
	CodeValidationError = "validation_error"
)

type Error struct {
	RequestID string `json:"request_id"`
	Status    int    `json:"status_code"`
	Code      string `json:"code"`
	Message   string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

func BadRequestError(message, requestID string) Error {
	return Error{
		Status:    http.StatusBadRequest,
		RequestID: requestID,
		Code:      CodeBadRequest,
		Message:   message,
	}
}

func ValidationFailedError(message, requestID string) Error {
	return Error{
		Status:    http.StatusBadRequest,
		RequestID: requestID,
		Code:      CodeValidationError,
		Message:   message,
	}
}
