package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"GrowEasy/dto"
)

type MLClient struct {
	BaseURL string
	Client  *http.Client
}

func NewMLClient() *MLClient {
	return &MLClient{
		BaseURL: "http://localhost:8000",
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Predict sends weather and soil data to ML model and returns prediction
func (c *MLClient) Predict(weather, soil map[string]interface{}) (*dto.MLServiceResponse, error) {
	reqBody := map[string]interface{}{
		"weather": weather,
		"soil":    soil,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.Client.Post(c.BaseURL+"/predict", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to call ML service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service returned status: %d", resp.StatusCode)
	}

	var result dto.MLServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode ML response: %w", err)
	}

	return &result, nil
}

// MLClientInterface defines the contract for ML prediction
type MLClientInterface interface {
	Predict(weather, soil map[string]interface{}) (*dto.MLServiceResponse, error)
}

// Ensure NewMLClient implements MLClientInterface
var _ MLClientInterface = (*MLClient)(nil)
