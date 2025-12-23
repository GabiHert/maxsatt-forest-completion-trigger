package usecase

import (
	"context"
	"math"
	"time"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
)

// CalculateWeatherMetrics defines the use case for computing historical weather statistics.
type CalculateWeatherMetrics interface {
	// Calculate computes statistical weather metrics from historical weather data.
	// It calculates averages, standard deviations, totals, and consecutive dry days.
	Calculate(ctx context.Context, weatherData []entity.WeatherData, targetDate time.Time) (*entity.WeatherMetrics, error)
}

// calculateWeatherMetrics implements the CalculateWeatherMetrics usecase.
type calculateWeatherMetrics struct {
}

// NewCalculateWeatherMetricsUseCase creates a new CalculateWeatherMetrics usecase implementation.
func NewCalculateWeatherMetricsUseCase() CalculateWeatherMetrics {
	return &calculateWeatherMetrics{}
}

// Calculate computes statistical weather metrics from historical weather data over the configured window period.
// It calculates:
// - Average temperature and standard deviation
// - Average humidity and standard deviation
// - Total precipitation
// - Maximum consecutive dry days (days with 0mm precipitation)
func (u *calculateWeatherMetrics) Calculate(ctx context.Context, weatherData []entity.WeatherData, targetDate time.Time) (*entity.WeatherMetrics, error) {
	logger.Debug(ctx, "Started", "targetDate", targetDate.Format("2006-01-02"), "dataPoints", len(weatherData))

	if len(weatherData) == 0 {
		return &entity.WeatherMetrics{}, nil
	}

	temperatures := make([]float64, 0, len(weatherData))
	humidities := make([]float64, 0, len(weatherData))
	precipitations := make([]float64, 0, len(weatherData))

	for _, data := range weatherData {
		temperatures = append(temperatures, data.Temperature)
		humidities = append(humidities, data.Humidity)
		precipitations = append(precipitations, data.Precipitation)
	}

	avgTemp := calculateAverage(temperatures)
	tempStdDev := calculateStandardDeviation(temperatures, avgTemp)
	avgHumidity := calculateAverage(humidities)
	humidityStdDev := calculateStandardDeviation(humidities, avgHumidity)
	totalPrecip := calculateSum(precipitations)
	dryDaysConsecutive := calculateMaxConsecutiveDryDays(precipitations)

	metrics := &entity.WeatherMetrics{
		AvgTemperature:     avgTemp,
		TempStdDev:         tempStdDev,
		AvgHumidity:        avgHumidity,
		HumidityStdDev:     humidityStdDev,
		TotalPrecipitation: totalPrecip,
		DryDaysConsecutive: dryDaysConsecutive,
	}

	logger.Debug(ctx, "Finished",
		"avgTemp", metrics.AvgTemperature,
		"tempStdDev", metrics.TempStdDev,
		"avgHumidity", metrics.AvgHumidity,
		"humidityStdDev", metrics.HumidityStdDev,
		"totalPrecip", metrics.TotalPrecipitation,
		"dryDaysConsecutive", metrics.DryDaysConsecutive,
	)
	return metrics, nil
}

// calculateAverage computes the mean of a slice of float64 values.
func calculateAverage(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}
	return roundToTwoDecimals(calculateSum(values) / float64(len(values)))
}

// calculateSum computes the sum of a slice of float64 values.
func calculateSum(values []float64) float64 {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum
}

// calculateStandardDeviation computes the standard deviation of a slice of float64 values.
// It uses the sample standard deviation formula: sqrt(sum((x - mean)^2) / (n - 1))
func calculateStandardDeviation(values []float64, mean float64) float64 {
	if len(values) <= 1 {
		return 0.0
	}

	variance := 0.0
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	variance /= float64(len(values) - 1)

	return roundToTwoDecimals(math.Sqrt(variance))
}

// calculateMaxConsecutiveDryDays finds the maximum number of consecutive days without precipitation.
// A dry day is defined as a day with 0.0mm precipitation.
func calculateMaxConsecutiveDryDays(precipitations []float64) int32 {
	var maxConsecutive int32
	var currentConsecutive int32

	for _, precip := range precipitations {
		if precip == 0.0 {
			currentConsecutive++
			if currentConsecutive > maxConsecutive {
				maxConsecutive = currentConsecutive
			}
		} else {
			currentConsecutive = 0
		}
	}

	return maxConsecutive
}

// roundToTwoDecimals rounds a float64 value to 2 decimal places.
// Uses standard rounding (0.5 rounds up).
func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}
