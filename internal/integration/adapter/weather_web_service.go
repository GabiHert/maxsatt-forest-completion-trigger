package adapter

import (
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/application/usecase"
)

// WeatherWebService defines the interface for fetching weather data from external APIs.
type WeatherWebService interface {
	usecase.FetchWeatherData
}
