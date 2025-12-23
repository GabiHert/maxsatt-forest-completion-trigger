package dto

import (
	"fmt"
	"time"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
)

// WeatherAPIResponse represents the response from the Open-Meteo Archive API.
// It contains latitude, longitude, daily weather data, and hourly weather data.
type WeatherAPIResponse struct {
	Daily     DailyWeatherData  `json:"daily"`
	Hourly    HourlyWeatherData `json:"hourly"`
	Latitude  float64           `json:"latitude"`
	Longitude float64           `json:"longitude"`
}

// DailyWeatherData contains the daily weather measurements.
type DailyWeatherData struct {
	Time          []string  `json:"time"`
	Temperature   []float64 `json:"temperature_2m_mean"`
	Precipitation []float64 `json:"precipitation_sum"`
}

// HourlyWeatherData contains the hourly weather measurements.
type HourlyWeatherData struct {
	Time     []string  `json:"time"`
	Humidity []float64 `json:"relative_humidity_2m"`
}

// ToEntity converts the Open-Meteo API response DTO to domain weather data entities.
// It calculates daily humidity from hourly data and combines all weather metrics.
func (w *WeatherAPIResponse) ToEntity() ([]entity.WeatherData, error) {
	if len(w.Daily.Time) == 0 {
		return nil, fmt.Errorf("no daily weather data in API response")
	}

	dailyHumidity := w.calculateDailyHumidity()

	weatherData := make([]entity.WeatherData, 0, len(w.Daily.Time))

	for i, dateStr := range w.Daily.Time {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse date %s: %w", dateStr, err)
		}

		temperature := 0.0
		if i < len(w.Daily.Temperature) {
			temperature = w.Daily.Temperature[i]
		}

		precipitation := 0.0
		if i < len(w.Daily.Precipitation) {
			precipitation = w.Daily.Precipitation[i]
		}

		humidity := 0.0
		if humidityVal, exists := dailyHumidity[dateStr]; exists {
			humidity = humidityVal
		}

		weatherData = append(weatherData, entity.WeatherData{
			Date:          date,
			Temperature:   temperature,
			Humidity:      humidity,
			Precipitation: precipitation,
		})
	}

	return weatherData, nil
}

// calculateDailyHumidity calculates daily mean humidity from hourly humidity data.
func (w *WeatherAPIResponse) calculateDailyHumidity() map[string]float64 {
	dailyHumidity := make(map[string]float64)
	dailyCounts := make(map[string]int)

	for i, timestampStr := range w.Hourly.Time {
		if i >= len(w.Hourly.Humidity) {
			break
		}

		if len(timestampStr) < 10 {
			continue
		}
		dateStr := timestampStr[:10]

		humidity := w.Hourly.Humidity[i]
		dailyHumidity[dateStr] += humidity
		dailyCounts[dateStr]++
	}

	for date := range dailyHumidity {
		if count := dailyCounts[date]; count > 0 {
			dailyHumidity[date] /= float64(count)
		}
	}

	return dailyHumidity
}
