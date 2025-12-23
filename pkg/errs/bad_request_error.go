package errs

import (
	"errors"
	"net/http"
)

type badRequestError struct {
	message      string
	description  string
	code         string
	errorDetails []ErrorDetails
	statusCode   int
}

func BadRequestError(description, code string, details ...ErrorDetails) BaseError {
	err := badRequestError{
		statusCode:   http.StatusBadRequest,
		message:      "bad request",
		description:  description,
		code:         code,
		errorDetails: details,
	}

	return err
}

func (b badRequestError) ToDto(transactionId string) Error {
	return Error{
		Error: ErrorData{
			Description:  b.description,
			Code:         b.code,
			Id:           transactionId,
			ErrorDetails: b.errorDetails,
		},
	}
}

func (b badRequestError) Error() string {
	return b.message
}

func (b badRequestError) Description() string {
	return b.description
}

func (b badRequestError) StatusCode() int {
	return b.statusCode
}

func (b badRequestError) SetStatusCode(statusCode int, description string) {
	b.statusCode = statusCode
	b.description = description
}

func (b badRequestError) Is(target error) bool {
	return errors.As(target, &b)
}

func (b badRequestError) Code() string {
	return b.code
}

func (b badRequestError) ErrorDetails() map[string]any {
	var details = make(map[string]any)
	for _, detail := range b.errorDetails {
		details["attribute"] = detail.Attribute
		details["messages"] = detail.Messages
	}
	return details
}

func (badRequestError) Retries() int {
	return 0
}

func (badRequestError) InternalNotify() bool {
	return true
}

func (badRequestError) Abort() bool {
	return true
}
