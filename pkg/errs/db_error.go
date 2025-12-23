package errs

import (
	"errors"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

type dbError struct {
	message      string
	description  string
	code         string
	errorDetails []ErrorDetails
	statusCode   int
}

func DbError(err error) BaseError {
	var dbErr dbError
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return NotFoundError(err.Error(), "NF-01404")
	case errors.Is(err, gorm.ErrCheckConstraintViolated) || strings.Contains(err.Error(), "constraint failed: UNIQUE constraint failed: "):
		dbErr = dbError{
			statusCode:   http.StatusBadRequest,
			message:      "bad request",
			description:  "Record already exists",
			code:         "BR-01400",
			errorDetails: nil,
		}
	default:
		return InternalServerError(err, "IS-01500", "Something went wrong, please try again later")
	}

	return dbErr
}

func (d dbError) ToDto(transactionId string) Error {
	return Error{
		Error: ErrorData{
			Description:  d.description,
			Code:         d.code,
			Id:           transactionId,
			ErrorDetails: d.errorDetails,
		},
	}
}

func (d dbError) Error() string {
	return d.message
}

func (d dbError) Description() string {
	return d.description
}

func (d dbError) StatusCode() int {
	return d.statusCode
}

func (d dbError) SetStatusCode(statusCode int, description string) {
	d.statusCode = statusCode
	d.description = description
}

func (d dbError) Is(target error) bool {
	return errors.As(target, &d)
}

func (d dbError) Code() string {
	return d.code
}

func (d dbError) ErrorDetails() map[string]any {
	var details = make(map[string]any)
	for _, detail := range d.errorDetails {
		details["attribute"] = detail.Attribute
		details["messages"] = detail.Messages
	}
	return details
}

func (dbError) Retries() int {
	return 5
}

func (dbError) InternalNotify() bool {
	return true
}

func (dbError) Abort() bool {
	return false
}
