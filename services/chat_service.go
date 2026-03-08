package services

import (
	"GrowEasy/config"
	"GrowEasy/models"
	"GrowEasy/services/integration"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

// GetOrCreateSessionID gets the current session ID for a user, creating one if none exists
func (s *ChatService) GetOrCreateSessionID(userID string) (string, error) {
	var latestMessage models.ChatMessage
	result := config.DB.Where("user_id = ?", userID).Order("created_at DESC").First(&latestMessage)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return "", result.Error
	}
	if result.RowsAffected > 0 && latestMessage.SessionID != "" {
		return latestMessage.SessionID, nil
	}
	// No messages or empty session_id, create new session
	return s.CreateNewSession(userID)
}

// CreateNewSession creates a new session ID for a user
func (s *ChatService) CreateNewSession(userID string) (string, error) {
	sessionID := uuid.New().String()
	return sessionID, nil
}

// SaveMessage saves a chat message to DB
func (s *ChatService) SaveMessage(userID string, message string, isUser bool) error {
	sessionID, err := s.GetOrCreateSessionID(userID)
	if err != nil {
		return err
	}
	chatMsg := models.ChatMessage{
		UserID:    &userID,
		SessionID: sessionID,
		Message:   message,
		IsUser:    isUser,
	}
	return config.DB.Create(&chatMsg).Error
}

// GetChatHistory retrieves chat history for the current session of a user
func (s *ChatService) GetChatHistory(userID string) ([]models.ChatMessage, error) {
	sessionID, err := s.GetOrCreateSessionID(userID)
	if err != nil {
		return nil, err
	}
	var messages []models.ChatMessage
	result := config.DB.
		Preload("User").
		Where("user_id = ? AND session_id = ?", userID, sessionID).
		Order("created_at ASC"). // Display Oldest to Newest
		Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}

// GetAllChatHistory retrieves all chat history for a user, grouped by session
func (s *ChatService) GetAllChatHistory(userID string) (map[string][]models.ChatMessage, error) {
	var messages []models.ChatMessage
	result := config.DB.
		Preload("User").
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&messages)
	if result.Error != nil {
		return nil, result.Error
	}

	// Group by session
	grouped := make(map[string][]models.ChatMessage)
	for _, msg := range messages {
		grouped[msg.SessionID] = append(grouped[msg.SessionID], msg)
	}
	return grouped, nil
}

// ResetSession clears the chat session for a user and starts a new one
func (s *ChatService) ResetSession(userID string) error {
	s.geminiClient.ClearSession(userID)

	_, err := s.CreateNewSession(userID)
	return err
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
