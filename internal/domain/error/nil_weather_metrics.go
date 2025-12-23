package error

import (
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
)

// NewNilWeatherMetricsError creates a new NilWeatherMetricsError.
func NewNilWeatherMetricsError() error {
	return errs.BadRequestError("weather metrics cannot be nil", "CLIMATE-01-0010", errs.ErrorDetails{
		Attribute: "weather_metrics",
		Messages:  []string{"INVALID_VALUE"},
	})
}
