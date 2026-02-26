package main

import (
	"oat431/big-shwein-api/internal/config"
	"oat431/big-shwein-api/internal/route"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func main() {
	config.LoadEnvConfig()
	app := fiber.New()
	route.SetupMainRoute(app)

	port := os.Getenv("PORT")
	err := app.Listen(":" + port)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
