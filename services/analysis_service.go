package service

import (
	"encoding/json"
	"fmt"

	"GrowEasy/config"
	"GrowEasy/dto"
	"GrowEasy/models"
	"GrowEasy/services/integration"

	"gorm.io/datatypes"
)

type AnalysisService struct {
	weatherClient *integration.WeatherClient
	soilClient    *integration.SoilClient
	mlClient      *integration.MLClient
}

func NewAnalysisService() *AnalysisService {
	return &AnalysisService{
		weatherClient: integration.NewWeatherClient(),
		soilClient:    integration.NewSoilClient(),
		mlClient:      integration.NewMLClient(),
	}
}

func (s *AnalysisService) Create(userID string, req dto.AnalysisResponse) error {

	soilJSON, _ := json.Marshal(req.SoilData)
	weatherJSON, _ := json.Marshal(req.WeatherData)
	predictionJSON, _ := json.Marshal(req.Predictions)

	analysis := models.Analysis{
		UserID:       &userID,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		LocationName: req.LocationName,
		SoilData:     datatypes.JSON(soilJSON),
		WeatherData:  datatypes.JSON(weatherJSON),
		Predictions:  datatypes.JSON(predictionJSON),
		AiResponse:   req.AiResponse,
	}

	return config.DB.Create(&analysis).Error
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
	}

	if err := config.DB.Create(&analysis).Error; err != nil {
		return nil, fmt.Errorf("failed to save analysis: %w", err)
	}

	return &analysis, nil
}
