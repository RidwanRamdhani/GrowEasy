package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type SoilClient struct {
	BaseURL string
	Client  *http.Client
}

func NewSoilClient() *SoilClient {
	return &SoilClient{
		BaseURL: "https://rest.isric.org/soilgrids/v2.0",
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *SoilClient) GetSoilData(lat, lng float64) (map[string]interface{}, error) {
	// Build the URL with query parameters
	endpoint := "/properties/query"

	// Create URL with proper query parameters
	baseURL, err := url.Parse(c.BaseURL + endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	params := url.Values{}
	params.Set("lon", fmt.Sprintf("%f", lng))
	params.Set("lat", fmt.Sprintf("%f", lat))
	params.Set("property", "nitrogen")
	params.Set("property", "phh2o")
	params.Set("depth", "0-5cm")
	params.Set("depth", "5-15cm")
	params.Set("depth", "15-30cm")
	params.Set("value", "mean")

	baseURL.RawQuery = params.Encode()

	resp, err := c.Client.Get(baseURL.String())
	if err != nil {
		return nil, fmt.Errorf("failed to fetch soil data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("soil API returned status: %d", resp.StatusCode)
	}

	// Decode raw JSON response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode soil response: %w", err)
	}

	return result, nil
}

// SoilClientInterface defines the contract for soil data fetching
type SoilClientInterface interface {
	GetSoilData(lat, lng float64) (map[string]interface{}, error)
}

// Ensure NewSoilClient implements SoilClientInterface
var _ SoilClientInterface = (*SoilClient)(nil)
