package error

import (
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
)

// NewEmptyDatasetError creates a new EmptyDatasetError.
func NewEmptyDatasetError() error {
	return errs.BadRequestError("delta dataset is empty", "CLIMATE-01-0004", errs.ErrorDetails{
		Attribute: "dataset",
		Messages:  []string{"EMPTY_DATASET"},
	})
}
