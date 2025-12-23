package entity

import "time"

// WeatherData represents historical weather information for a specific location and date.
// It contains temperature, humidity, and precipitation data.
type WeatherData struct {
	Date          time.Time
	Temperature   float64
	Humidity      float64
	Precipitation float64
}
