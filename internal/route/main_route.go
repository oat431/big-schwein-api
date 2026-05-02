package route

import (
	"oat431/big-shwein-api/internal/bootstrap"
	"oat431/big-shwein-api/internal/config"
	"oat431/big-shwein-api/internal/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func init() {
	log.Info("Initializing main route...")
}
func SetupMainRoute(app *fiber.App, apiContainer *bootstrap.APIContainer) {
	app.Use(middleware.GlobalMiddleware)
	app.Use(cors.New(config.InitCorsConfig()))

	api := app.Group("/api")
	v1 := api.Group("/v1")

	RegisterHealthCheckRoute(v1)
	RegisterAuthRoute(v1, apiContainer.AuthController)
}
