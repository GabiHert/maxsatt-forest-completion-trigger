package errs

import (
	"errors"
)

type serverError struct {
	message      string
	description  string
	code         string
	errorDetails []ErrorDetails
	statusCode   int
}

func ServerError(description, code string, status int) BaseError {
	err := serverError{
		statusCode:   status,
		message:      "server error occurred",
		description:  description,
		code:         code,
		errorDetails: nil,
	}

	return err
}

func (s serverError) ToDto(transactionId string) Error {
	return Error{
		Error: ErrorData{
			Description:  s.description,
			Code:         s.code,
			Id:           transactionId,
			ErrorDetails: s.errorDetails,
		},
	}
}

func (s serverError) Error() string {
	return s.message
}

func (s serverError) Description() string {
	return s.description
}

func (s serverError) StatusCode() int {
	return s.statusCode
}

func (s serverError) SetStatusCode(statusCode int, description string) {
	s.statusCode = statusCode
	s.description = description
}

func (s serverError) Is(target error) bool {
	return errors.As(target, &s)
}

func (s serverError) Code() string {
	return s.code
}

func (s serverError) ErrorDetails() map[string]any {
	var details = make(map[string]any)
	for _, detail := range s.errorDetails {
		details["attribute"] = detail.Attribute
		details["messages"] = detail.Messages
	}
	return details
}

func (serverError) Retries() int {
	return 5
}

func (serverError) InternalNotify() bool {
	return true
}

func (serverError) Abort() bool {
	return false
}
