package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	AuthID uuid.UUID `json:"auth_id"`
	jwt.RegisteredClaims
}

func getJWTSecret() ([]byte, error) {
	secret := os.Getenv("SECRET")
	if secret == "" {
		return nil, errors.New("missing jwt secret")
	}
	return []byte(secret), nil
}

func GenerateAccessToken(authID uuid.UUID) (string, error) {
	jwtSecret, err := getJWTSecret()
	if err != nil {
		return "", err
	}

	claims := JWTClaims{
		AuthID: authID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "big-shwein",
			Subject:   authID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken() (string, error) {
	return GenerateUUID(), nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	jwtSecret, err := getJWTSecret()
	if err != nil {
		return nil, err
	}

	return jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}
