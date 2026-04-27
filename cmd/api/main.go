package main

import (
	"oat431/big-shwein-api/internal/bootstrap"
	"oat431/big-shwein-api/internal/config"
	"oat431/big-shwein-api/internal/route"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

func main() {
	config.LoadEnvConfig()
	db := config.StartDatabase()
	defer db.Close()

	apiContainer, err := bootstrap.NewAPIContainer(db)
	if err != nil {
		log.Fatal("Failed to initialize API container: ", err)
	}

	app := fiber.New()
	route.SetupMainRoute(app, apiContainer)

	port := os.Getenv("PORT")
	err = app.Listen(":" + port)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
