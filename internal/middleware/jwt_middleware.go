package middleware

import (
	"strings"

	"oat431/big-shwein-api/pkg/common"
	"oat431/big-shwein-api/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

func JWTMiddleware(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.StatusUnauthorized,
				ErrorCode: "UNAUTHORIZED",
				Message:   "missing authorization header",
			},
		})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.StatusUnauthorized,
				ErrorCode: "UNAUTHORIZED",
				Message:   "invalid authorization header format",
			},
		})
	}

	tokenString := parts[1]
	token, err := utils.ValidateToken(tokenString)
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.StatusUnauthorized,
				ErrorCode: "UNAUTHORIZED",
				Message:   "invalid or expired token",
			},
		})
	}

	claims, ok := token.Claims.(*utils.JWTClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.StatusUnauthorized,
				ErrorCode: "UNAUTHORIZED",
				Message:   "invalid token claims",
			},
		})
	}

	c.Locals("auth_id", claims.AuthID)

	return c.Next()
}
