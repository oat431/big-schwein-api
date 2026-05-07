package config

import (
	"github.com/gofiber/fiber/v3/log"
	"github.com/joho/godotenv"
)

// LoadEnvConfig loads environment variables from the .env.development file.
func LoadEnvConfig() {
	failedToLoadEnv := godotenv.Load(".env.development")
	if failedToLoadEnv != nil {
		log.Fatal("Failed to load environment variables from .env.development file")
	}
}
