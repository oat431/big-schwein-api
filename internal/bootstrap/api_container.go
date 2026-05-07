package bootstrap

import (
	"oat431/big-shwein-api/internal/config"
	"oat431/big-shwein-api/internal/controller"
	"oat431/big-shwein-api/internal/repository"
	"oat431/big-shwein-api/internal/service"

	"github.com/gofiber/fiber/v3/log"
	"github.com/jmoiron/sqlx"
)

// APIContainer holds all controller instances for dependency injection.
type APIContainer struct {
	AuthController *controller.AuthController
}

// NewAPIContainer wires all dependencies and returns a fully initialized APIContainer.
func NewAPIContainer(db *sqlx.DB) (*APIContainer, error) {
	log.Info("registering repositories")
	authRepository := repository.NewAuthRepository(db)
	refreshTokenRepository := repository.NewRefreshTokenRepository(db)
	emailVerifyTokenRepository := repository.NewEmailVerifyTokenRepository(db)

	log.Info("registering services")
	smtpService := service.NewSMTPService(config.GetSMTPConfig())
	tokenConfig := config.GetTokenConfig()
	authService, err := service.NewAuthService(
		authRepository,
		refreshTokenRepository,
		emailVerifyTokenRepository,
		smtpService,
		tokenConfig,
	)
	if err != nil {
		return nil, err
	}

	log.Info("registering controllers")
	authController, err := controller.NewAuthController(authService)
	if err != nil {
		return nil, err
	}

	log.Info("all auth dependencies registered")
	return &APIContainer{AuthController: authController}, nil
}
