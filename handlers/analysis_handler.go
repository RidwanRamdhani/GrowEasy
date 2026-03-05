package handlers

import (
	"net/http"

	"GrowEasy/dto"
	service "GrowEasy/services"

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

// GetWeather fetches weather data from Open-Meteo API
func (h *AnalysisHandler) GetWeather(c *gin.Context) {
	var req dto.AnalysisRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	weather, err := h.service.GetWeather(req.Latitude, req.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, weather)
}

// GetSoil fetches soil data from SoilGrids API
func (h *AnalysisHandler) GetSoil(c *gin.Context) {
	var req dto.AnalysisRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	soil, err := h.service.GetSoil(req.Latitude, req.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, soil)
}

// GetPredict performs ML prediction and saves to DB
func (h *AnalysisHandler) GetPredict(c *gin.Context) {
	var req dto.AnalysisRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(string)

	analysis, err := h.service.GetPredict(userID, req.Latitude, req.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analysis)
}

// GetHistory fetches all analysis history for the authenticated user
func (h *AnalysisHandler) GetHistory(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	history, err := h.service.GetAnalysisHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  history,
		"count": len(history),
	})
}

