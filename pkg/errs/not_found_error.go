package errs

import (
	"errors"
	"net/http"
)

type notFoundError struct {
	message      string
	description  string
	code         string
	errorDetails []ErrorDetails
	statusCode   int
}

func NotFoundError(message string, code string) BaseError {
	return notFoundError{
		statusCode:   http.StatusNotFound,
		message:      "not found",
		description:  message,
		code:         code,
		errorDetails: nil,
	}
}

func IsNotFoundError(err error) bool {
	return notFoundError{}.Is(err)
}

func (n notFoundError) ToDto(transactionId string) Error {
	return Error{
		Error: ErrorData{
			Description:  n.description,
			Code:         n.code,
			Id:           transactionId,
			ErrorDetails: n.errorDetails,
		},
	}
}

func (n notFoundError) Error() string {
	return n.message
}

func (n notFoundError) Description() string {
	return n.description
}

func (n notFoundError) StatusCode() int {
	return n.statusCode
}

func (n notFoundError) SetStatusCode(statusCode int, description string) {
	n.statusCode = statusCode
	n.description = description
}

func (n notFoundError) Is(target error) bool {
	return errors.As(target, &n)
}

func (n notFoundError) Code() string {
	return n.code
}

func (n notFoundError) ErrorDetails() map[string]any {
	var details = make(map[string]any)
	for _, detail := range n.errorDetails {
		details["attribute"] = detail.Attribute
		details["messages"] = detail.Messages
	}
	return details
}

func (notFoundError) Retries() int {
	return 0
}

func (notFoundError) InternalNotify() bool {
	return true
}

func (notFoundError) Abort() bool {
	return true
}
