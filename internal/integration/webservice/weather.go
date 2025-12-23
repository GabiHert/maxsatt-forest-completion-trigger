package webservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/domain/entity"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/infra/properties"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/adapter"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/internal/integration/webservice/dto"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/errs"
	http2 "github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/http"
	"github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/logger"
	redishelper "github.com/GabiHert/maxsatt-forest-completion-trigger/pkg/redis"
)

// weatherWebService implements the WeatherAPIClient interface.
type weatherWebService struct {
	httpClient http2.Client
	redis      redishelper.RedisHelper
}

// NewWeather creates a new WeatherAPIClient implementation.
func NewWeather(httpClient http2.Client, redis redishelper.RedisHelper) adapter.WeatherWebService {
	return &weatherWebService{
		httpClient: httpClient,
		redis:      redis,
	}
}

// FetchWeather retrieves historical weather data from Open-Meteo Archive API.
// It fetches data for a historical window period based on the HISTORICAL_CLIMATE_DAYS configuration.
// The window is calculated as: startDate = targetDate - HISTORICAL_CLIMATE_DAYS, endDate = targetDate.
func (a *weatherWebService) FetchWeather(ctx context.Context, lat, lon float64, date time.Time) ([]entity.WeatherData, error) {
	historicalDays := properties.Properties().Services.HistoricalClimateDays
	startDate := date.AddDate(0, 0, -historicalDays)
	endDate := date

	cacheKey := fmt.Sprintf("weather:%.4f:%.4f:%s:%s", lat, lon, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	var cachedData []entity.WeatherData
	if err := a.redis.Get(ctx, cacheKey, &cachedData); err == nil && len(cachedData) > 0 {
		return cachedData, nil
	}

	params := url.Values{}
	params.Set("latitude", fmt.Sprintf("%.4f", lat))
	params.Set("longitude", fmt.Sprintf("%.4f", lon))
	params.Set("start_date", startDate.Format("2006-01-02"))
	params.Set("end_date", endDate.Format("2006-01-02"))
	params.Set("daily", "temperature_2m_mean,precipitation_sum")
	params.Set("hourly", "relative_humidity_2m")

	apiURL := fmt.Sprintf("%s?%s", properties.Properties().Services.WeatherApiUrl+"/v1/archive", params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, errs.InternalServerError(err, "WEA-00500", "Failed to create weather data request")
	}

	resp, err := a.httpClient.Do(ctx, req)
	if err != nil {
		return nil, errs.InternalServerError(err, "WEA-01500", "Failed to fetch weather data")
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		return nil, errs.ClientError("Invalid weather data request", "WEA-00400", resp.StatusCode)
	} else if resp.StatusCode == 429 {
		return nil, errs.ServerError("Weather API rate limit exceeded", "WEA-03500", resp.StatusCode)
	} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		return nil, errs.ClientError("Error while trying to get weather data", "WEA-00400", resp.StatusCode)
	} else if resp.StatusCode >= 500 {
		return nil, errs.ServerError("weather API returned server error", "WEA-02500", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(ctx, err, "Failed to read weather API response body", "url", apiURL)
		return nil, errs.InternalServerError(err, "WEA-03500", "Failed to read weather data response body")
	}

	var apiResponse dto.WeatherAPIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		logger.Error(ctx, err, "Failed to parse weather API JSON response", "url", apiURL, "response_body", string(body))
		return nil, errs.InternalServerError(err, "WEA-04500", "Failed to parse weather data response")
	}

	weatherData, err := apiResponse.ToEntity()
	if err != nil {
		logger.Error(ctx, err, "Failed to convert weather API response to entity", "url", apiURL)
		return nil, errs.InternalServerError(err, "WEA-05500", "Failed to convert weather data")
	}

	_ = a.redis.Put(ctx, cacheKey, weatherData, 24*time.Hour)

	return weatherData, nil
}
