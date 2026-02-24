package handler

import (
	"net/http"

	"GrowEasy/dto"
	"GrowEasy/service"

	"github.com/gin-gonic/gin"
)

type AnalysisHandler struct {
	service *service.AnalysisService
}

func NewAnalysisHandler() *AnalysisHandler {
	return &AnalysisHandler{
		service: service.NewAnalysisService(),
	}
}

func (h *AnalysisHandler) Create(c *gin.Context) {

	var req dto.CreateAnalysisRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil user_id dari JWT
	userID := c.MustGet("user_id").(string)

	if err := h.service.Create(userID, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Analysis saved"})
}