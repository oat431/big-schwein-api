package utils

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims defines the custom claims for access tokens.
type JWTClaims struct {
	AuthID uuid.UUID `json:"auth_id"`
	jwt.RegisteredClaims
}

var (
	jwtSecret     []byte
	jwtSecretOnce sync.Once
	jwtSecretErr  error
)

func loadJWTSecret() {
	secret := os.Getenv("SECRET")
	if secret == "" {
		jwtSecretErr = errors.New("missing jwt secret")
		return
	}
	jwtSecret = []byte(secret)
}

func getJWTSecret() ([]byte, error) {
	jwtSecretOnce.Do(loadJWTSecret)
	if jwtSecretErr != nil {
		return nil, jwtSecretErr
	}
	return jwtSecret, nil
}

// GenerateAccessToken creates a signed JWT access token for the given auth ID.
func GenerateAccessToken(authID uuid.UUID, expiry time.Duration) (string, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return "", fmt.Errorf("generate access token: %w", err)
	}

	now := time.Now()
	claims := JWTClaims{
		AuthID: authID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "big-shwein",
			Subject:   authID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// GenerateRefreshToken creates a new random refresh token string.
func GenerateRefreshToken() (string, error) {
	return GenerateUUID(), nil
}

// ValidateToken parses and validates a JWT token string.
func ValidateToken(tokenString string) (*jwt.Token, error) {
	secret, err := getJWTSecret()
	if err != nil {
		return nil, fmt.Errorf("validate token: %w", err)
	}

	return jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}
