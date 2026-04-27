package route

import (
	"oat431/big-shwein-api/internal/controller"
	"oat431/big-shwein-api/internal/middleware"
	"oat431/big-shwein-api/internal/payload/request"

	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoute(router fiber.Router, controller *controller.AuthController) {
	route := router.Group("/auth")

	route.Post("/register",
		middleware.Validate[request.RegisterRequest],
		controller.RegisterNewUser,
	)

	route.Post("/login",
		middleware.Validate[request.LoginRequest],
		controller.Login,
	)

	route.Post("/revoke", controller.RevokeAccess)

	route.Get("/detail",
		middleware.JWTMiddleware,
		controller.GetUserDetails,
	)

	route.Get("/verify-email", controller.VerifyEmail)
}
