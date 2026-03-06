package services

import (
	"GrowEasy/config"
	"GrowEasy/models"
	"GrowEasy/services/integration"
	"fmt"
)

type ChatService struct {
	analysisService *AnalysisService
	geminiClient    *integration.GeminiClient
}

func NewChatService() *ChatService {
	geminiClient, _ := integration.NewGeminiClient() // Handle error as needed
	return &ChatService{
		analysisService: NewAnalysisService(),
		geminiClient:    geminiClient,
	}
}

// SaveMessage saves a chat message to DB
func (s *ChatService) SaveMessage(userID string, message string, isUser bool) error {
	chatMsg := models.ChatMessage{
		UserID:  &userID,
		Message: message,
		IsUser:  isUser,
	}
	return config.DB.Create(&chatMsg).Error
}

// GetChatHistory retrieves chat history for a user
func (s *ChatService) GetChatHistory(userID string) ([]models.ChatMessage, error) {
	var messages []models.ChatMessage
	result := config.DB.
		Where("user_id = ?", userID).
		Order("created_at ASC"). // Display Oldest to Newest
		Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}

// ResetSession clears the chat session for a user
func (s *ChatService) ResetSession(userID string) error {
	s.geminiClient.ClearSession(userID)
	return nil
}

// Chat handles user message with context from latest analysis using Gemini session
func (s *ChatService) Chat(userID, message string) (string, error) {
	analysis, err := s.analysisService.GetLatestAnalysis(userID)
	if err != nil {
		return "", fmt.Errorf("no analysis found: %w", err)
	}

	// Fetch base knowledge from latest analysis
	context := analysis.AiResponse

	if err := s.SaveMessage(userID, message, true); err != nil {
		return "", fmt.Errorf("failed to save user message: %w", err)
	}

	// Get response from Gemini session
	response, err := s.geminiClient.ChatWithSession(userID, context, message)
	if err != nil {
		return "", err
	}

	if err := s.SaveMessage(userID, response, false); err != nil {
		return "", fmt.Errorf("failed to save AI response: %w", err)
	}

	return response, nil
}
