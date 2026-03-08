package handlers

import (
	"GrowEasy/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	service *services.ChatService
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{
		service: services.NewChatService(),
	}
}

// Chat processes user messages and returns AI responses based on latest analysis
func (h *ChatHandler) Chat(c *gin.Context) {
	var req struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(string)

	response, err := h.service.Chat(userID, req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}

// GetHistory retrieves and returns the user's chat message history
func (h *ChatHandler) GetHistory(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	all := c.Query("all") == "true"

	if all {
		history, err := h.service.GetAllChatHistory(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Count total messages across all sessions
		totalCount := 0
		for _, msgs := range history {
			totalCount += len(msgs)
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  history,
			"count": totalCount,
		})
	} else {
		history, err := h.service.GetChatHistory(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data":  history,
			"count": len(history),
		})
	}
}

// Reset clears the user's chat session to refresh context
func (h *ChatHandler) Reset(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	err := h.service.ResetSession(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chat session reset successfully"})
}
