package handler

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

func (h *AnalysisHandler) Create(c *gin.Context) {

	var req dto.AnalysisResponse

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
