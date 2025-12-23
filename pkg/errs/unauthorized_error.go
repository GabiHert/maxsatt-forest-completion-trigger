package errs

import (
	"errors"
	"net/http"
)

type unauthorizedError struct {
	message      string
	description  string
	code         string
	errorDetails []ErrorDetails
	statusCode   int
}

func UnauthorizedError(message string, code string) BaseError {
	return unauthorizedError{
		statusCode:   http.StatusUnauthorized,
		message:      message,
		description:  "Unauthorized",
		code:         code,
		errorDetails: nil,
	}
}

func (u unauthorizedError) ToDto(transactionId string) Error {
	return Error{
		Error: ErrorData{
			Description:  u.description,
			Code:         u.code,
			Id:           transactionId,
			ErrorDetails: u.errorDetails,
		},
	}
}

func (u unauthorizedError) Error() string {
	return u.message
}

func (u unauthorizedError) Description() string {
	return u.description
}

func (u unauthorizedError) StatusCode() int {
	return u.statusCode
}

func (u unauthorizedError) SetStatusCode(statusCode int, description string) {
	u.statusCode = statusCode
	u.description = description
}

func (u unauthorizedError) Is(target error) bool {
	return errors.As(target, &u)
}

func (u unauthorizedError) Code() string {
	return u.code
}

func (u unauthorizedError) ErrorDetails() map[string]any {
	var details = make(map[string]any)
	for _, detail := range u.errorDetails {
		details["attribute"] = detail.Attribute
		details["messages"] = detail.Messages
	}
	return details
}

func (unauthorizedError) Retries() int {
	return 0
}

func (unauthorizedError) InternalNotify() bool {
	return true
}

func (unauthorizedError) Abort() bool {
	return true
}
