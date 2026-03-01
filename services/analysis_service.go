package service

import (
	"encoding/json"
	"fmt"

	"GrowEasy/config"
	"GrowEasy/models"
	"GrowEasy/services/integration"

	"gorm.io/datatypes"
)

type AnalysisService struct {
	weatherClient *integration.WeatherClient
	soilClient    *integration.SoilClient
	mlClient      *integration.MLClient
	geminiClient  *integration.GeminiClient
}

func NewAnalysisService() *AnalysisService {
	geminiClient, err := integration.NewGeminiClient()
	if err != nil {
		fmt.Printf("Warning: Failed to initialize Gemini client: %v\n", err)
	}

	return &AnalysisService{
		weatherClient: integration.NewWeatherClient(),
		soilClient:    integration.NewSoilClient(),
		mlClient:      integration.NewMLClient(),
		geminiClient:  geminiClient,
	}
}

// GetWeather fetches weather data from Open-Meteo API
func (s *AnalysisService) GetWeather(lat, lon float64) (map[string]interface{}, error) {
	weatherData, err := s.weatherClient.GetWeather(lat, lon)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather data: %w", err)
	}
	return weatherData, nil
}

// GetSoil fetches soil data from SoilGrids API
func (s *AnalysisService) GetSoil(lat, lon float64) (map[string]interface{}, error) {
	soilData, err := s.soilClient.GetSoilData(lat, lon)
	if err != nil {
		return nil, fmt.Errorf("failed to get soil data: %w", err)
	}
	return soilData, nil
}

// GetPredict fetches weather, soil, calls ML, and saves to DB
func (s *AnalysisService) GetPredict(userID string, lat, lon float64) (*models.Analysis, error) {
	// Fetch weather data
	weatherData, err := s.weatherClient.GetWeather(lat, lon)
	if err != nil {
		return nil, fmt.Errorf("failed to get weather data: %w", err)
	}

	// Fetch soil data
	soilData, err := s.soilClient.GetSoilData(lat, lon)
	if err != nil {
		return nil, fmt.Errorf("failed to get soil data: %w", err)
	}

	// Call ML model
	prediction, err := s.mlClient.Predict(weatherData, soilData)
	if err != nil {
		return nil, fmt.Errorf("failed to get ML prediction: %w", err)
	}

	// Generate summary with Gemini
	var aiResponse string
	if s.geminiClient != nil {
		summary, err := s.geminiClient.GenerateSummary(
			prediction.PredictionClass,
			prediction.Probability,
			prediction.Top3,
			soilData,
			weatherData,
		)
		if err != nil {
			// Log error but don't fail the whole request
			fmt.Printf("Warning: Failed to generate Gemini summary: %v\n", err)
			aiResponse = fmt.Sprintf("Prediction: %s (%.2f%%)", prediction.PredictionClass, prediction.Probability*100)
		} else {
			aiResponse = summary
		}
	} else {
		// Fallback if Gemini client is not available
		aiResponse = fmt.Sprintf("Prediction: %s (%.2f%%)", prediction.PredictionClass, prediction.Probability*100)
	}

	// Save to DB
	soilJSON, _ := json.Marshal(soilData)
	weatherJSON, _ := json.Marshal(weatherData)
	predictionJSON, _ := json.Marshal(prediction)

	analysis := models.Analysis{
		UserID:      &userID,
		Latitude:    lat,
		Longitude:   lon,
		SoilData:    datatypes.JSON(soilJSON),
		WeatherData: datatypes.JSON(weatherJSON),
		Predictions: datatypes.JSON(predictionJSON),
		AiResponse:  aiResponse,
	}

	if err := config.DB.Create(&analysis).Error; err != nil {
		return nil, fmt.Errorf("failed to save analysis: %w", err)
	}

	// Load the user
	user := &models.User{}
	if err := config.DB.First(user, "id = ?", *analysis.UserID).Error; err != nil {
		return nil, fmt.Errorf("failed to load user: %w", err)
	}
	analysis.User = *user

	return &analysis, nil
}

// UpdateAiResponse updates the AI summary for an existing analysis
func (s *AnalysisService) UpdateAiResponse(analysisID string, aiResponse string) error {
	return config.DB.Model(&models.Analysis{}).Where("id = ?", analysisID).Update("ai_response", aiResponse).Error
}

// GenerateSummary generates a summary using Gemini API
func (s *AnalysisService) GenerateSummary(
	predictionClass string,
	probability float64,
	top3 map[string]float64,
	soilData map[string]interface{},
	weatherData map[string]interface{},
) (string, error) {
	if s.geminiClient == nil {
		return "", fmt.Errorf("Gemini client is not initialized")
	}
	return s.geminiClient.GenerateSummary(predictionClass, probability, top3, soilData, weatherData)
}
