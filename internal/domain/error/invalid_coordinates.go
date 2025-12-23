package error

import (
	"fmt"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
)

// NewInvalidCoordinatesError creates a new InvalidCoordinatesError.
func NewInvalidCoordinatesError(lat, lon float64) error {
	return errs.BadRequestError(
		fmt.Sprintf("invalid coordinates: latitude=%f, longitude=%f (must be within -90 to 90 and -180 to 180)", lat, lon),
		"CLIMATE-01-0005",
		errs.ErrorDetails{
			Attribute: "coordinates",
			Messages:  []string{"INVALID_RANGE"},
		},
	)
}
