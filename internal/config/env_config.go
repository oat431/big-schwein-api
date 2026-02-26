package config

import (
	"github.com/gofiber/fiber/v3/log"
	"github.com/joho/godotenv"
)

func init() {
	log.Info("Loading environment variables...")
}

func LoadEnvConfig() {
	failedToLoadEnv := godotenv.Load(".env.development")
	if failedToLoadEnv != nil {
		log.Fatal("Failed to load environment variables from .env.development file")
	}
}
