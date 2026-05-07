package controller

import (
	"errors"

	"oat431/big-shwein-api/internal/httputil"
	"oat431/big-shwein-api/internal/payload/request"
	"oat431/big-shwein-api/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

// AuthController handles authentication HTTP endpoints.
type AuthController struct {
	service service.AuthService
}

// NewAuthController creates a new AuthController with the given service.
func NewAuthController(service service.AuthService) (*AuthController, error) {
	if service == nil {
		return nil, errors.New("auth controller: nil service")
	}
	return &AuthController{service: service}, nil
}

// RegisterNewUser handles POST /auth/register.
func (auth *AuthController) RegisterNewUser(c fiber.Ctx) error {
	req := c.Locals("payload").(*request.RegisterRequest)
	authDTO, err := auth.service.Register(c.Context(), *req)
	if err != nil {
		return httputil.ErrorResponse(c, fiber.StatusBadRequest, "REGISTER-01", err.Error())
	}

	return httputil.SuccessResponse(c, fiber.StatusCreated, authDTO)
}

// Login handles POST /auth/login.
func (auth *AuthController) Login(c fiber.Ctx) error {
	req := c.Locals("payload").(*request.LoginRequest)
	tokenDTO, err := auth.service.Login(c.Context(), *req)
	if err != nil {
		return httputil.ErrorResponse(c, fiber.StatusBadRequest, "LOGIN-01", "invalid username or password")
	}

	return httputil.SuccessResponse(c, fiber.StatusOK, tokenDTO)
}

// RevokeAccess handles POST /auth/revoke.
func (auth *AuthController) RevokeAccess(c fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.Bind().JSON(&body); err != nil {
		return httputil.ErrorResponse(c, fiber.StatusBadRequest, "REVOKE-01", "invalid request body")
	}

	if err := auth.service.RevokeAccess(c.Context(), body.RefreshToken); err != nil {
		return httputil.ErrorResponse(c, fiber.StatusInternalServerError, "REVOKE-02", err.Error())
	}

	message := "access revoked"
	return httputil.SuccessResponse(c, fiber.StatusOK, &message)
}

// GetUserDetails handles GET /auth/detail.
func (auth *AuthController) GetUserDetails(c fiber.Ctx) error {
	authID := c.Locals("auth_id").(uuid.UUID)
	authDTO, err := auth.service.GetUserDetails(c.Context(), authID)
	if err != nil {
		return httputil.ErrorResponse(c, fiber.StatusInternalServerError, "DETAIL-01", err.Error())
	}

	return httputil.SuccessResponse(c, fiber.StatusOK, authDTO)
}

// VerifyEmail handles GET /auth/verify-email.
func (auth *AuthController) VerifyEmail(c fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return httputil.ErrorResponse(c, fiber.StatusBadRequest, "VERIFY-01", "missing verification token")
	}

	if err := auth.service.VerifyEmail(c.Context(), token); err != nil {
		return httputil.ErrorResponse(c, fiber.StatusBadRequest, "VERIFY-02", err.Error())
	}

	message := "email verified successfully"
	return httputil.SuccessResponse(c, fiber.StatusOK, &message)
}
