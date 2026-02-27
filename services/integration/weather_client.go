package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type WeatherClient struct {
	BaseURL string
	Client  *http.Client
}

func NewWeatherClient() *WeatherClient {
	return &WeatherClient{
		BaseURL: "https://archive-api.open-meteo.com/v1",
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *WeatherClient) GetWeather(lat, lon float64) (map[string]interface{}, error) {
	// Calculate date range: 30 days ago to today
	endDate := time.Now().UTC().Format("2006-01-02")
	startDate := time.Now().UTC().AddDate(0, 0, -30).Format("2006-01-02")

	// Build URL with proper query parameters
	baseURL, err := url.Parse(c.BaseURL + "/archive")
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	params := url.Values{}
	params.Set("latitude", fmt.Sprintf("%f", lat))
	params.Set("longitude", fmt.Sprintf("%f", lon))
	params.Set("start_date", startDate)
	params.Set("end_date", endDate)
	params.Set("daily", "temperature_2m_mean,precipitation_sum")
	params.Set("hourly", "relative_humidity_2m")
	params.Set("timezone", "auto")

	baseURL.RawQuery = params.Encode()

	resp, err := c.Client.Get(baseURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned status: %d", resp.StatusCode)
	}

	// Return raw JSON as map[string]interface{}
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode weather response: %w", err)
	}

	return result, nil
}

// WeatherClientInterface defines the contract for weather data fetching
type WeatherClientInterface interface {
	GetWeather(lat, lon float64) (map[string]interface{}, error)
}

// Ensure NewWeatherClient implements WeatherClientInterface
var _ WeatherClientInterface = (*WeatherClient)(nil)
