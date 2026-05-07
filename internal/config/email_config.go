package config

import (
	"os"
	"strconv"
	"time"
)

// SMTPConfig holds SMTP server configuration.
type SMTPConfig struct {
	SMTPHost           string
	SMTPPort           int
	SMTPUser           string
	SMTPPassword       string
	CodeExpiration     time.Duration
	VerifyEmailBaseURL string
}

// GetSMTPConfig reads SMTP settings from environment variables.
func GetSMTPConfig() *SMTPConfig {
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		port = 587
	}

	verifyBaseURL := os.Getenv("VERIFY_EMAIL_BASE_URL")
	if verifyBaseURL == "" {
		verifyBaseURL = "http://localhost:3000/verify-email"
	}

	return &SMTPConfig{
		SMTPHost:           os.Getenv("SMTP_HOST"),
		SMTPPort:           port,
		SMTPUser:           os.Getenv("SMTP_USER"),
		SMTPPassword:       os.Getenv("SMTP_PASS"),
		CodeExpiration:     time.Minute,
		VerifyEmailBaseURL: verifyBaseURL,
	}
}
