package config

import (
	"os"
	"time"
)

// TokenConfig holds token expiration durations.
type TokenConfig struct {
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	VerifyTokenExpiry  time.Duration
}

// GetTokenConfig reads token expiration settings from environment variables
// with sensible defaults.
func GetTokenConfig() *TokenConfig {
	return &TokenConfig{
		AccessTokenExpiry:  parseDuration("ACCESS_TOKEN_EXPIRY", time.Hour),
		RefreshTokenExpiry: parseDuration("REFRESH_TOKEN_EXPIRY", 7*24*time.Hour),
		VerifyTokenExpiry:  parseDuration("VERIFY_TOKEN_EXPIRY", 24*time.Hour),
	}
}

func parseDuration(envKey string, defaultVal time.Duration) time.Duration {
	val := os.Getenv(envKey)
	if val == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return defaultVal
	}
	return d
}
