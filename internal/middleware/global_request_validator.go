package middleware

import (
	"fmt"
	"strings"

	"oat431/big-shwein-api/pkg/common"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate = validator.New()

func Validate[T any](c fiber.Ctx) error {
	payload := new(T)

	if err := c.Bind().Body(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.StatusBadRequest,
				ErrorCode: "BAD_REQUEST",
				Message:   "invalid request body format",
			},
		})
	}

	if err := validate.Struct(payload); err != nil {
		var errorMessages []string
		for _, validationErr := range err.(validator.ValidationErrors) {
			msg := fmt.Sprintf("Field '%s' failed on tag '%s'", validationErr.Field(), validationErr.Tag())
			errorMessages = append(errorMessages, msg)
		}

		return c.Status(fiber.StatusBadRequest).JSON(common.ResponseDTO[any]{
			Status: common.ERROR,
			Error: &common.ResponseDTOError{
				HttpCode:  fiber.StatusBadRequest,
				ErrorCode: "VALIDATION_ERROR",
				Message:   strings.Join(errorMessages, ", "),
			},
		})
	}

	c.Locals("payload", payload)

	return c.Next()
}
