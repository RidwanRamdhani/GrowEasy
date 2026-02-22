package main

import (
	"log"
	"net/http"

	"GrowEasy/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to database
	if err := config.ConnectDatabase(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer config.CloseDatabase()

	// Auto migrate models
	if err := config.AutoMigrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, Gin!",
		})
	})

	r.Run(":8080")
}
