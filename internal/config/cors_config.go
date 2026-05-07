package config

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v3/middleware/cors"
)

// InitCorsConfig creates CORS middleware configuration from environment variables.
func InitCorsConfig() cors.Config {
	origins := os.Getenv("CORS_ORIGINS")
	if origins == "" {
		origins = "http://localhost:3000,http://localhost:5173"
	}

	return cors.Config{
		AllowOrigins:     strings.Split(origins, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}
}
