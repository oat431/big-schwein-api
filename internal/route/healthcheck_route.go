package route

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
)

func init() {
	log.Info("healthcheck api initialized: /api/v1/health")
}
func RegisterHealthCheckRoute(router fiber.Router) {
	route := router.Group("/health")

	route.Get("/", healthcheck.New())
}
