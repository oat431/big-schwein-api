package httputil

import "github.com/gofiber/fiber/v3"

// ResponseDTOStatus represents the status of an API response.
type ResponseDTOStatus string

const (
	SUCCESS ResponseDTOStatus = "SUCCESS"
	FAIL    ResponseDTOStatus = "FAIL"
	ERROR   ResponseDTOStatus = "ERROR"
)

// ResponseDTOError holds error details for an API response.
type ResponseDTOError struct {
	HttpCode  int    `json:"http_code"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

// ResponseDTO is a generic API response wrapper.
type ResponseDTO[T any] struct {
	Data   *T                `json:"data"`
	Status ResponseDTOStatus `json:"status"`
	Error  *ResponseDTOError `json:"error"`
}

// ErrorResponse sends a standardized error JSON response.
func ErrorResponse(c fiber.Ctx, httpCode int, errorCode, message string) error {
	return c.Status(httpCode).JSON(ResponseDTO[any]{
		Status: ERROR,
		Error: &ResponseDTOError{
			HttpCode:  httpCode,
			ErrorCode: errorCode,
			Message:   message,
		},
	})
}

// SuccessResponse sends a standardized success JSON response.
func SuccessResponse[T any](c fiber.Ctx, httpCode int, data *T) error {
	return c.Status(httpCode).JSON(ResponseDTO[T]{
		Status: SUCCESS,
		Data:   data,
	})
}
