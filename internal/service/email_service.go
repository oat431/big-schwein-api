package service

import (
	"fmt"
	"net/smtp"

	"oat431/big-shwein-api/internal/config"
)

type SMTPService struct {
	config *config.Config
}

func NewSMTPService(cfg *config.Config) *SMTPService {
	return &SMTPService{config: cfg}
}

func (s *SMTPService) SendVerificationEmail(to, token string) error {
	subject := "Subject: Big Schwein Account Verification\n"
	body := fmt.Sprintf(`
		If you are receiving this email, your account has been successfully created.
		Please verify your email address by clicking this link: http://localhost:3000/verify-email?token=%s
		Thank you for joining Big Schwein!
	`, token)

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
