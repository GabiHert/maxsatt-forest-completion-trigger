package errors

import (
	"errors"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"

	"github.com/go-playground/validator/v10"
)

func ValidationError(err error, code string) error {
	var validationErrs validator.ValidationErrors
	var errorDetails []errs.ErrorDetails
	if ok := errors.As(err, &validationErrs); ok {
		for _, e := range validationErrs {
			switch {
			case e.Tag() == "not_nil" || e.Tag() == "required":
				errorDetails = append(errorDetails, errs.ErrorDetails{
					Attribute: e.Field(),
					Messages:  []string{"REQUIRED_ATTRIBUTE_MISSING"},
				})
			case e.Tag() == "uuid4" || e.Tag() == "uuid":
				errorDetails = append(errorDetails, errs.ErrorDetails{
					Attribute: e.Field(),
					Messages:  []string{"INVALID_UUID_FORMAT"},
				})
			case e.Tag() == "gt":
				errorDetails = append(errorDetails, errs.ErrorDetails{
					Attribute: e.Field(),
					Messages:  []string{"INVALID_VALUE"},
				})
			default:
				errorDetails = append(errorDetails, errs.ErrorDetails{
					Attribute: e.Field(),
					Messages:  []string{"UNKNOWN_VALIDATION_ERROR"},
				})
			}
		}
	} else {
		errorDetails = append(errorDetails, errs.ErrorDetails{
			Attribute: "",
			Messages:  []string{"UNKNOWN_VALIDATION_ERROR"},
		})
	}
	return errs.BadRequestError("Bad request", code, errorDetails...)
}
