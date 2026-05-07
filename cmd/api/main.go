package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"oat431/big-shwein-api/internal/bootstrap"
	"oat431/big-shwein-api/internal/config"
	"oat431/big-shwein-api/internal/route"

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

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Fatal("Failed to start server: ", err)
		}
	}()

	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Info("Server exited gracefully")
}
