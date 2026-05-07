package middleware

import (
	"strings"

	"oat431/big-shwein-api/internal/httputil"
	"oat431/big-shwein-api/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

// JWTMiddleware validates the Authorization header and extracts auth claims.
func JWTMiddleware(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return httputil.ErrorResponse(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "missing authorization header")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return httputil.ErrorResponse(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "invalid authorization header format")
	}

	tokenString := parts[1]
	token, err := utils.ValidateToken(tokenString)
	if err != nil || !token.Valid {
		return httputil.ErrorResponse(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "invalid or expired token")
	}

	claims, ok := token.Claims.(*utils.JWTClaims)
	if !ok {
		return httputil.ErrorResponse(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "invalid token claims")
	}

	c.Locals("auth_id", claims.AuthID)

	return c.Next()
}
