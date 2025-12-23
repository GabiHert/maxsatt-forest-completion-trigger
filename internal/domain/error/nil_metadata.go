package error

import (
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
)

// NewNilMetadataError creates a new NilMetadataError.
func NewNilMetadataError() error {
	return errs.BadRequestError("metadata cannot be nil", "CLIMATE-01-0009", errs.ErrorDetails{
		Attribute: "metadata",
		Messages:  []string{"INVALID_VALUE"},
	})
}
