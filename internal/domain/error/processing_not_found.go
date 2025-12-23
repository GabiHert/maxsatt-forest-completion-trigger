package error

import (
	"fmt"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
)

// NewProcessingNotFoundError creates a new ProcessingNotFoundError.
func NewProcessingNotFoundError(processingID string) error {
	return errs.NotFoundError(
		fmt.Sprintf("processing data not found for ID: %s", processingID),
		"CLIMATE-01-0002",
	)
}
