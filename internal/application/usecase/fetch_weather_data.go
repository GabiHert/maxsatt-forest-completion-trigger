package usecase

import (
	"context"
	"time"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
)

// FetchWeatherData defines the use case for retrieving weather data with caching.
type FetchWeatherData interface {
	// FetchWeather fetches weather data for the given location and time range.
	// It checks the cache first and falls back to the API if necessary.
	FetchWeather(ctx context.Context, lat, lon float64, date time.Time) ([]entity.WeatherData, error)
}
