package services

import (
	"fmt"

	"GrowEasy/config"
	"GrowEasy/models"
)

// GetAnalysisHistory returns all analyses for a given user, newest first
func (s *AnalysisService) GetAnalysisHistory(userID string) ([]models.Analysis, error) {
	var analyses []models.Analysis

	result := config.DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&analyses)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to fetch analysis history: %w", result.Error)
	}

	return analyses, nil
}
