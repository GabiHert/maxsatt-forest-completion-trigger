package errs

import (
	"errors"
	"net/http"
)

type forbiddenError struct {
	message      string
	description  string
	code         string
	errorDetails []ErrorDetails
	statusCode   int
}

func ForbiddenError(message string, code string) BaseError {
	return forbiddenError{
		statusCode:   http.StatusForbidden,
		message:      message,
		description:  "Forbidden.",
		code:         code,
		errorDetails: nil,
	}
}

func (f forbiddenError) ToDto(transactionId string) Error {
	return Error{
		Error: ErrorData{
			Description:  f.description,
			Code:         f.code,
			Id:           transactionId,
			ErrorDetails: f.errorDetails,
		},
	}
}

func (f forbiddenError) Error() string {
	return f.message
}

func (f forbiddenError) Description() string {
	return f.description
}

func (f forbiddenError) StatusCode() int {
	return f.statusCode
}

func (f forbiddenError) SetStatusCode(statusCode int, description string) {
	f.statusCode = statusCode
	f.description = description
}

func (f forbiddenError) Is(target error) bool {
	return errors.As(target, &f)
}

func (f forbiddenError) Code() string {
	return f.code
}

func (f forbiddenError) ErrorDetails() map[string]any {
	var details = make(map[string]any)
	for _, detail := range f.errorDetails {
		details["attribute"] = detail.Attribute
		details["messages"] = detail.Messages
	}
	return details
}

func (f forbiddenError) Retries() int {
	return 0
}

func (f forbiddenError) InternalNotify() bool {
	return false
}

func (f forbiddenError) Abort() bool {
	return true
}
