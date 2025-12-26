package error

import (
	"fmt"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
)

// NewNotificationFailedError creates a new error when a notification fails to be published
// for a specific forest.
func NewNotificationFailedError(forestId string, err error) error {
	return errs.InternalServerError(
		err,
		"COMPLETION-01-0002",
		fmt.Sprintf("failed to publish notification for forest: %s", forestId),
	)
}
