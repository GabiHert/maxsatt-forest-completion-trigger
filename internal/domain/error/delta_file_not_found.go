package error

import (
	"fmt"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
)

// NewDeltaFileNotFoundError creates a new error when DELTA file is not found for a processing ID.
func NewDeltaFileNotFoundError(processingID string) error {
	return errs.NotFoundError(
		fmt.Sprintf("DELTA file not found for processing ID: %s", processingID),
		"CLIMATE-01-0005",
	)
}
