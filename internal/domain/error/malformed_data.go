package error

import (
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
)

// NewMalformedDataError creates a new MalformedDataError.
func NewMalformedDataError() error {
	return errs.BadRequestError("malformed data - JSON parsing failed", "CLIMATE-01-0006", errs.ErrorDetails{
		Attribute: "data",
		Messages:  []string{"INVALID_FORMAT"},
	})
}
