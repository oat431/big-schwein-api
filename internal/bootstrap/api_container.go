package bootstrap

import (
	"oat431/big-shwein-api/internal/config"
	"oat431/big-shwein-api/internal/controller"
	"oat431/big-shwein-api/internal/repository"
	"oat431/big-shwein-api/internal/service"

	"github.com/gofiber/fiber/v3/log"
	"github.com/jmoiron/sqlx"
)

type APIContainer struct {
	AuthController *controller.AuthController
}

func NewAPIContainer(db *sqlx.DB) (*APIContainer, error) {
	log.Info("registering repositories")
	authRepository := repository.NewAuthRepository(db)
	refreshTokenRepository := repository.NewRefreshTokenRepository(db)
	emailVerifyTokenRepository := repository.NewEmailVerifyTokenRepository(db)

	log.Info("registering services")
	smtpService := service.NewSMTPService(config.GetEmailConfig())
	authService, err := service.NewAuthService(authRepository, refreshTokenRepository, emailVerifyTokenRepository, smtpService)
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
