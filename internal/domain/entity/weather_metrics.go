package entity

// WeatherMetrics represents historical weather data statistics for a specific location and time period.
// These metrics are calculated from historical weather data over a configurable window period.
// Statistical calculations provide insights into climate patterns and variability.
type WeatherMetrics struct {
	// AvgTemperature is the average temperature across the historical window period (degrees Celsius)
	AvgTemperature float64

	// TempStdDev is the standard deviation of temperatures across the historical window period (degrees Celsius)
	TempStdDev float64

	// AvgHumidity is the average relative humidity across the historical window period (percentage)
	AvgHumidity float64

	// HumidityStdDev is the standard deviation of humidity across the historical window period (percentage)
	HumidityStdDev float64

	// TotalPrecipitation is the sum of all precipitation across the historical window period (millimeters)
	TotalPrecipitation float64

	// DryDaysConsecutive is the maximum number of consecutive days without precipitation (days with 0mm)
	DryDaysConsecutive int32
}
