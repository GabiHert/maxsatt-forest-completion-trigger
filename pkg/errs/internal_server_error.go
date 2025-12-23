package errs

import (
	"errors"
	"net/http"
)

type internalServerError struct {
	message      string
	description  string
	code         string
	errorDetails []ErrorDetails
	statusCode   int
}

func InternalServerError(err error, code, message string) BaseError {
	return internalServerError{
		statusCode:  http.StatusInternalServerError,
		message:     "internal server",
		description: message,
		code:        code,
		errorDetails: []ErrorDetails{
			{
				Attribute: "UNKNOWN_ERROR",
				Messages:  []string{err.Error()},
			},
		},
	}
}

func (i internalServerError) ToDto(transactionId string) Error {
	return Error{
		Error: ErrorData{
			Description:  i.description,
			Code:         i.code,
			Id:           transactionId,
			ErrorDetails: i.errorDetails,
		},
	}
}

func (i internalServerError) Error() string {
	return i.message
}

func (i internalServerError) Description() string {
	return i.description
}

func (i internalServerError) StatusCode() int {
	return i.statusCode
}

func (i internalServerError) SetStatusCode(statusCode int, description string) {
	i.statusCode = statusCode
	i.description = description
}

func (i internalServerError) Is(target error) bool {
	return errors.As(target, &i)
}

func (i internalServerError) Code() string {
	return i.code
}

func (i internalServerError) ErrorDetails() map[string]any {
	var details = make(map[string]any)
	for _, detail := range i.errorDetails {
		details["attribute"] = detail.Attribute
		details["messages"] = detail.Messages
	}
	return details
}

func (internalServerError) Retries() int {
	return 5
}

func (internalServerError) InternalNotify() bool {
	return true
}

func (internalServerError) Abort() bool {
	return false
}
