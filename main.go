package main

import (
	"log"
	"net/http"

	"GrowEasy/config"
	handler "GrowEasy/handlers"
	"GrowEasy/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.ConnectDatabase(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer config.CloseDatabase()

	if err := config.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, Gin!",
		})
	})

	authHandler := handler.NewAuthHandler()
	analysisHandler := handler.NewAnalysisHandler()

	auth := r.Group("/api/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			c.JSON(http.StatusOK, gin.H{
				"message": "Authorized",
				"user_id": userID,
			})
		})

		// Weather data endpoint
		api.POST("/weather", analysisHandler.GetWeather)

		// Soil data endpoint
		api.POST("/soil", analysisHandler.GetSoil)

		// Unified analysis: weather + soil + ML prediction
		api.POST("/predict", analysisHandler.GetPredict)
	}

	r.Run(":8080")
}
