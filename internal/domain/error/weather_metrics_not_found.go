package error

import (
	"fmt"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
)

// NewWeatherMetricsNotFoundError creates a new WeatherMetricsNotFoundError.
func NewWeatherMetricsNotFoundError(date string) error {
	return errs.BadRequestError(
		fmt.Sprintf("no weather metrics found for date: %s", date),
		"CLIMATE-01-0008",
		errs.ErrorDetails{
			Attribute: "weather_metrics",
			Messages:  []string{"METRICS_NOT_FOUND"},
		},
	)
}
