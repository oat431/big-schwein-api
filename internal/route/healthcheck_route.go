package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
)

// RegisterHealthCheckRoute sets up the /health endpoint.
func RegisterHealthCheckRoute(router fiber.Router) {
	route := router.Group("/health")

	route.Get("/", healthcheck.New())
}
