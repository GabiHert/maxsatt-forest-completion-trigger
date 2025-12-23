package error

import (
	"fmt"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
)

// NewMissingImageIDError creates a new error when DELTA file is missing the image_id field.
func NewMissingImageIDError(processingID string) error {
	return errs.BadRequestError(
		fmt.Sprintf("DELTA file for processing ID %s is missing image_id", processingID),
		"CLIMATE-01-0006",
		errs.ErrorDetails{
			Attribute: "image_id",
			Messages:  []string{"MISSING_IMAGE_ID"},
		},
	)
}
