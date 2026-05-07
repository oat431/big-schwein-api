package service

import (
	"fmt"
	"net/smtp"

	"oat431/big-shwein-api/internal/config"
)

// SMTPService sends emails via SMTP.
type SMTPService struct {
	config *config.SMTPConfig
}

// NewSMTPService creates a new SMTPService with the given configuration.
func NewSMTPService(cfg *config.SMTPConfig) *SMTPService {
	return &SMTPService{config: cfg}
}

// SendVerificationEmail sends an account verification email to the given address.
func (s *SMTPService) SendVerificationEmail(to, token string) error {
	subject := "Subject: Big Schwein Account Verification\n"
	body := fmt.Sprintf(`
		If you are receiving this email, your account has been successfully created.
		Please verify your email address by clicking this link: %s?token=%s
		Thank you for joining Big Schwein!
	`, s.config.VerifyEmailBaseURL, token)

	message := []byte(subject + "\n" + body)

	auth := smtp.PlainAuth(
		"",
		s.config.SMTPUser,
		s.config.SMTPPassword,
		s.config.SMTPHost,
	)

	return smtp.SendMail(
		fmt.Sprintf("%s:%d", s.config.SMTPHost, s.config.SMTPPort),
		auth,
		s.config.SMTPUser,
		[]string{to},
		message,
	)
}
