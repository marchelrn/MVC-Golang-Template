package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/smtp"
)

func GenerateToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

type SMTPConfig struct {
	Host     string
	Port     string
	Email    string
	Password string
	AppPort  int
}

func SendEmail(toEmail, token string, cfg SMTPConfig) error {
	if cfg.Email == "" || cfg.Password == "" {
		return fmt.Errorf("SMTP not configured: set SMTP_EMAIL and SMTP_PASSWORD in .env")
	}

	verifyLink := fmt.Sprintf("http://localhost:%d/verify-email?token=%s", cfg.AppPort, token)

	subject := "Email Verification - Mini Jira"
	body := fmt.Sprintf(`Hello!

Please click the following link to verify your email address:

%s

This link will expire in 24 hours.

If you did not register for Mini Jira, please ignore this email.`, verifyLink)

	msg := fmt.Sprintf(
		"From: Mini Jira <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"UTF-8\"\r\n\r\n%s",
		cfg.Email, toEmail, subject, body,
	)

	auth := smtp.PlainAuth("", cfg.Email, cfg.Password, cfg.Host)
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	return smtp.SendMail(addr, auth, cfg.Email, []string{toEmail}, []byte(msg))
}
