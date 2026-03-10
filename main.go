package main

import (
	"log"
	"net/http"
	"time"

	"GrowEasy/config"
	handler "GrowEasy/handlers"
	"GrowEasy/middleware"

	"github.com/gin-contrib/cors"
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

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, Gin!",
		})
	})

	authHandler := handler.NewAuthHandler()
	analysisHandler := handler.NewAnalysisHandler()
	chatHandler := handler.NewChatHandler()

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

		// Unified analysis: weather + soil + ML prediction + Gemini Summary
		api.POST("/predict", analysisHandler.GetPredict)

		// Fetch all analysis history for authenticated user
		api.GET("/history", analysisHandler.GetHistory)

		// Chat with Gemini using latest analysis context
		api.POST("/chat", chatHandler.Chat)

		// Get chat history
		api.GET("/chat/history", chatHandler.GetHistory)

		// Reset chat session
		api.POST("/chat/reset", chatHandler.Reset)
	}

	r.Run(":8080")
}
