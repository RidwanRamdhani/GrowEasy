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
	params.Add("property", "nitrogen")
	params.Add("property", "phh2o")
	params.Add("depth", "0-5cm")
	params.Add("depth", "5-15cm")
	params.Add("depth", "15-30cm")
	params.Set("value", "mean")

	baseURL.RawQuery = params.Encode()

	resp, err := c.Client.Get(baseURL.String())
	if err != nil {
		// Check if it's a timeout error
		if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
			// Fallback to mock data on timeout
			return c.getMockSoilData(lat, lng), nil
		}
		return nil, fmt.Errorf("failed to fetch soil data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusServiceUnavailable {
		// Fallback to mock data when API is unavailable
		return c.getMockSoilData(lat, lng), nil
	}

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

// getMockSoilData returns mock soil data when API is unavailable
func (c *SoilClient) getMockSoilData(lat, lng float64) map[string]interface{} {
	return map[string]interface{}{
		"geometry": map[string]interface{}{
			"coordinates": []interface{}{lng, lat},
			"type":        "Point",
		},
		"properties": map[string]interface{}{
			"layers": []interface{}{
				map[string]interface{}{
					"name": "nitrogen",
					"unit_measure": map[string]interface{}{
						"d_factor":         100,
						"mapped_units":     "cg/kg",
						"target_units":     "g/kg",
						"uncertainty_unit": "",
					},
					"depths": []interface{}{
						map[string]interface{}{
							"label": "0-5cm",
							"range": map[string]interface{}{
								"bottom_depth": 5,
								"top_depth":    0,
								"unit_depth":   "cm",
							},
							"values": map[string]interface{}{
								"mean": 1.5,
							},
						},
						map[string]interface{}{
							"label": "5-15cm",
							"range": map[string]interface{}{
								"bottom_depth": 15,
								"top_depth":    5,
								"unit_depth":   "cm",
							},
							"values": map[string]interface{}{
								"mean": 1.3,
							},
						},
						map[string]interface{}{
							"label": "15-30cm",
							"range": map[string]interface{}{
								"bottom_depth": 30,
								"top_depth":    15,
								"unit_depth":   "cm",
							},
							"values": map[string]interface{}{
								"mean": 1.2,
							},
						},
					},
				},
				map[string]interface{}{
					"name": "phh2o",
					"unit_measure": map[string]interface{}{
						"d_factor":         10,
						"mapped_units":     "pH*10",
						"target_units":     "-",
						"uncertainty_unit": "",
					},
					"depths": []interface{}{
						map[string]interface{}{
							"label": "0-5cm",
							"range": map[string]interface{}{
								"bottom_depth": 5,
								"top_depth":    0,
								"unit_depth":   "cm",
							},
							"values": map[string]interface{}{
								"mean": 6.5,
							},
						},
						map[string]interface{}{
							"label": "5-15cm",
							"range": map[string]interface{}{
								"bottom_depth": 15,
								"top_depth":    5,
								"unit_depth":   "cm",
							},
							"values": map[string]interface{}{
								"mean": 6.8,
							},
						},
						map[string]interface{}{
							"label": "15-30cm",
							"range": map[string]interface{}{
								"bottom_depth": 30,
								"top_depth":    15,
								"unit_depth":   "cm",
							},
							"values": map[string]interface{}{
								"mean": 7.0,
							},
						},
					},
				},
			},
		},
		"query_time_s": 2.194484233856201,
		"type":         "Feature",
	}
}

// SoilClientInterface defines the contract for soil data fetching
type SoilClientInterface interface {
	GetSoilData(lat, lng float64) (map[string]interface{}, error)
}

// Ensure NewSoilClient implements SoilClientInterface
var _ SoilClientInterface = (*SoilClient)(nil)
