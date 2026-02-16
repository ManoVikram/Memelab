package routes

import (
	"net/http"

	"github.com/ManoVikram/AI-Meme-Generator/backend/api/handlers"
	"github.com/ManoVikram/AI-Meme-Generator/backend/api/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine, services *services.Services) {
	// GET request for API health check (no auth required)
	server.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "fitness-coach-api",
		})
	})

	// GET method to generate meme image with a caption based on a topic
	server.POST("/api/generate", handlers.GenerateMemeHandler(services))
}
