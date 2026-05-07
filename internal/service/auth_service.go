package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"oat431/big-shwein-api/internal/config"
	"oat431/big-shwein-api/internal/model"
	"oat431/big-shwein-api/internal/payload/request"
	"oat431/big-shwein-api/internal/payload/response"
	"oat431/big-shwein-api/pkg/utils"

	"github.com/gofiber/fiber/v3/log"
	"github.com/google/uuid"
)

// AuthStore defines the data access contract for auth entities.
type AuthStore interface {
	Register(ctx context.Context, request request.RegisterRequest) (*model.Auth, error)
	GetAuthByUsername(ctx context.Context, username string) (*model.Auth, error)
	GetAuthByID(ctx context.Context, id uuid.UUID) (*model.Auth, error)
	GetAuthByEmail(ctx context.Context, email string) (*model.Auth, error)
	MarkAsVerified(ctx context.Context, authID uuid.UUID) error
}

// RefreshTokenStore defines the data access contract for refresh tokens.
type RefreshTokenStore interface {
	Save(ctx context.Context, refreshToken model.RefreshToken) error
	Revoke(ctx context.Context, token string) error
	GetByToken(ctx context.Context, token string) (*model.RefreshToken, error)
}

// EmailVerifyTokenStore defines the data access contract for email verification tokens.
type EmailVerifyTokenStore interface {
	Save(ctx context.Context, token model.EmailVerifyToken) error
	FindByToken(ctx context.Context, token string) (*model.EmailVerifyToken, error)
	DeleteByAuthID(ctx context.Context, authID uuid.UUID) error
}

// EmailSender defines the contract for sending emails.
type EmailSender interface {
	SendVerificationEmail(to, token string) error
}

// AuthService defines the business operations for authentication.
type AuthService interface {
	Register(ctx context.Context, request request.RegisterRequest) (*response.AuthResponse, error)
	Login(ctx context.Context, request request.LoginRequest) (*response.JWTResponse, error)
	RevokeAccess(ctx context.Context, refreshToken string) error
	GetUserDetails(ctx context.Context, authID uuid.UUID) (*response.AuthResponse, error)
	VerifyEmail(ctx context.Context, token string) error
}

type authService struct {
	repo                 AuthStore
	refreshTokenRepo     RefreshTokenStore
	emailVerifyTokenRepo EmailVerifyTokenStore
	emailSender          EmailSender
	tokenConfig          *config.TokenConfig
}

// NewAuthService creates a new AuthService with all required dependencies.
func NewAuthService(
	repo AuthStore,
	refreshTokenRepo RefreshTokenStore,
	emailVerifyTokenRepo EmailVerifyTokenStore,
	emailSender EmailSender,
	tokenConfig *config.TokenConfig,
) (AuthService, error) {
	if repo == nil {
		return nil, errors.New("auth service: nil auth repository")
	}
	if refreshTokenRepo == nil {
		return nil, errors.New("auth service: nil refresh token repository")
	}
	if emailVerifyTokenRepo == nil {
		return nil, errors.New("auth service: nil email verify token repository")
	}
	if emailSender == nil {
		return nil, errors.New("auth service: nil email sender")
	}
	if tokenConfig == nil {
		return nil, errors.New("auth service: nil token config")
	}

	return &authService{
		repo:                 repo,
		refreshTokenRepo:     refreshTokenRepo,
		emailVerifyTokenRepo: emailVerifyTokenRepo,
		emailSender:          emailSender,
		tokenConfig:          tokenConfig,
	}, nil
}

func (s *authService) Register(ctx context.Context, request request.RegisterRequest) (*response.AuthResponse, error) {
	encryptedPassword, err := utils.EncryptPassword(request.Password)
	if err != nil {
		return nil, fmt.Errorf("encrypt password: %w", err)
	}
	request.Password = encryptedPassword
	auth, err := s.repo.Register(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	tokenStr, err := utils.GenerateVerifyToken()
	if err != nil {
		return nil, fmt.Errorf("generate verification token: %w", err)
	}

	now := time.Now()
	emailVerifyToken := model.EmailVerifyToken{
		BaseEntity: model.BaseEntity{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		AuthID:    auth.ID,
		Token:     tokenStr,
		ExpiresAt: now.Add(s.tokenConfig.VerifyTokenExpiry),
	}

	if err := s.emailVerifyTokenRepo.Save(ctx, emailVerifyToken); err != nil {
		return nil, fmt.Errorf("save email verify token: %w", err)
	}

	if err := s.emailSender.SendVerificationEmail(auth.Email, tokenStr); err != nil {
		log.Error("failed to send verification email: ", err.Error())
	}

	return &response.AuthResponse{
		ID:         auth.ID.String(),
		Username:   auth.Username,
		Email:      auth.Email,
		IsVerified: auth.IsVerified,
	}, nil
}

func (s *authService) Login(ctx context.Context, request request.LoginRequest) (*response.JWTResponse, error) {
	auth, err := s.repo.GetAuthByUsername(ctx, request.Username)
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	if err := utils.ComparePassword(auth.Password, request.Password); err != nil {
		return nil, fmt.Errorf("login: invalid credentials: %w", err)
	}

	accessToken, err := utils.GenerateAccessToken(auth.ID, s.tokenConfig.AccessTokenExpiry)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshTokenStr, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	now := time.Now()
	refreshToken := model.RefreshToken{
		BaseEntity: model.BaseEntity{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		AuthID:    auth.ID,
		Token:     refreshTokenStr,
		ExpiresAt: now.Add(s.tokenConfig.RefreshTokenExpiry),
		Revoked:   false,
	}

	if err := s.refreshTokenRepo.Save(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	return &response.JWTResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
	}, nil
}

func (s *authService) RevokeAccess(ctx context.Context, refreshToken string) error {
	if err := s.refreshTokenRepo.Revoke(ctx, refreshToken); err != nil {
		return fmt.Errorf("revoke access: %w", err)
	}
	return nil
}

func (s *authService) GetUserDetails(ctx context.Context, authID uuid.UUID) (*response.AuthResponse, error) {
	auth, err := s.repo.GetAuthByID(ctx, authID)
	if err != nil {
		return nil, fmt.Errorf("get user details: %w", err)
	}

	return &response.AuthResponse{
		ID:         auth.ID.String(),
		Username:   auth.Username,
		Email:      auth.Email,
		IsVerified: auth.IsVerified,
	}, nil
}

func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	emailVerifyToken, err := s.emailVerifyTokenRepo.FindByToken(ctx, token)
	if err != nil {
		return errors.New("invalid or expired verification token")
	}

	if time.Now().After(emailVerifyToken.ExpiresAt) {
		return errors.New("verification token has expired")
	}

	if err := s.repo.MarkAsVerified(ctx, emailVerifyToken.AuthID); err != nil {
		return fmt.Errorf("mark as verified: %w", err)
	}

	if err := s.emailVerifyTokenRepo.DeleteByAuthID(ctx, emailVerifyToken.AuthID); err != nil {
		log.Error("failed to delete email verify token: ", err.Error())
	}

	return nil
}
