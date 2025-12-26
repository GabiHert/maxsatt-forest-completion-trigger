package error

import (
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
)

// NewNoCompletedForestsError creates a new error when no forests are found
// with all processings completed and pending notification.
func NewNoCompletedForestsError() error {
	return errs.NotFoundError(
		"no forests found with all processings completed and pending notification",
		"COMPLETION-01-0001",
	)
}
