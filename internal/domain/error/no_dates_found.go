package error

import (
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
)

// NewNoDatesFoundError creates a new NoDatesFoundError.
func NewNoDatesFoundError() error {
	return errs.BadRequestError("no dates found in delta dataset", "CLIMATE-01-0007", errs.ErrorDetails{
		Attribute: "dataset",
		Messages:  []string{"NO_DATES_FOUND"},
	})
}
