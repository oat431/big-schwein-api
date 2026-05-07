package middleware

import (
	"fmt"
	"strings"

	"oat431/big-shwein-api/internal/httputil"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate = validator.New()

// Validate is a generic middleware that binds and validates the request body
// against the given type T, storing the result in Locals("payload").
func Validate[T any](c fiber.Ctx) error {
	payload := new(T)

	if err := c.Bind().Body(payload); err != nil {
		return httputil.ErrorResponse(c, fiber.StatusBadRequest, "BAD_REQUEST", "invalid request body format")
	}

	if err := validate.Struct(payload); err != nil {
		var errorMessages []string
		for _, validationErr := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on tag '%s'", validationErr.Field(), validationErr.Tag())
			errorMessages = append(errorMessages, msg)
		}

		return httputil.ErrorResponse(c, fiber.StatusBadRequest, "VALIDATION_ERROR", strings.Join(errorMessages, ", "))
	}

	c.Locals("payload", payload)

	return c.Next()
}
