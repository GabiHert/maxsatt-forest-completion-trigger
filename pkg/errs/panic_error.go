package errs

import (
	"errors"
	"net/http"
)

type panicError struct {
	message      string
	description  string
	code         string
	errorDetails []ErrorDetails
	statusCode   int
}

func PanicError(message string) BaseError {
	err := panicError{
		statusCode:   http.StatusInternalServerError,
		message:      message,
		description:  "Something went wrong, please try again in a few seconds",
		code:         "IS-00500",
		errorDetails: nil,
	}

	return err
}

func (p panicError) ToDto(transactionId string) Error {
	return Error{
		Error: ErrorData{
			Description:  p.description,
			Code:         p.code,
			Id:           transactionId,
			ErrorDetails: p.errorDetails,
		},
	}
}

func (p panicError) Error() string {
	return p.message
}

func (p panicError) Description() string {
	return p.description
}

func (p panicError) StatusCode() int {
	return p.statusCode
}

func (p panicError) SetStatusCode(statusCode int, description string) {
	p.statusCode = statusCode
	p.description = description
}

func (p panicError) Retries() int {
	return 0
}

func (p panicError) InternalNotify() bool {
	return true
}

func (panicError) Abort() bool {
	return true
}

func (p panicError) Is(target error) bool {
	return errors.As(target, &p)
}

func (p panicError) Code() string {
	return p.code
}

func (p panicError) ErrorDetails() map[string]any {
	var details = make(map[string]any)
	for _, detail := range p.errorDetails {
		details["attribute"] = detail.Attribute
		details["messages"] = detail.Messages
	}
	return details
}
