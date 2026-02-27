package service

import (
	"encoding/json"

	"GrowEasy/config"
	"GrowEasy/dto"
	"GrowEasy/models"

	"gorm.io/datatypes"
)

type AnalysisService struct{}

func NewAnalysisService() *AnalysisService {
	return &AnalysisService{}
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
