package service

import (
	"fmt"
	"log"
	"mini_jira/config"
	"mini_jira/contract"
	"mini_jira/models"
	errs "mini_jira/pkg/error"
	"mini_jira/utils"
	"time"
)

type EmailVerificationService struct {
	verificationRepo contract.EmailVerificationRepository
	userRepo         contract.UserRepository
}

func ImplEmailVerificationService(
	verificationRepo contract.EmailVerificationRepository,
	userRepo contract.UserRepository,
) contract.EmailVerificationService {
	return &EmailVerificationService{
		verificationRepo: verificationRepo,
		userRepo:         userRepo,
	}
}

func (s *EmailVerificationService) VerifyEmail(token string) error {
	if token == "" {
		return errs.BadRequest("token is required")
	}

	verification, err := s.verificationRepo.GetByToken(token)
	if err != nil {
		return errs.BadRequest("invalid or expired token")
	}

	// Check expired
	if time.Now().After(verification.ExpiredAt) {
		s.verificationRepo.Delete(verification.ID)
		return errs.BadRequest("token has expired, please register again or request a new verification email")
	}

	// Update user status to active
	if err := s.userRepo.UpdateUserStatus(verification.UserID, "active"); err != nil {
		return errs.InternalServerError("failed to activate user")
	}

	// Delete used token
	s.verificationRepo.Delete(verification.ID)

	return nil
}

func (s *EmailVerificationService) ResendVerification(email string) error {
	if email == "" {
		return errs.BadRequest("email is required")
	}

	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return errs.NotFound("email not found")
	}

	if user.Status == "active" {
		return errs.BadRequest("email is already verified")
	}

	// Delete old tokens for this user
	s.verificationRepo.DeleteByUserID(user.Id)

	// Generate new token
	token, err := utils.GenerateToken()
	if err != nil {
		return errs.InternalServerError("failed to generate verification token")
	}

	verification := &models.EmailsVerification{
		UserID:    user.Id,
		Email:     user.Email,
		Token:     token,
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.verificationRepo.Create(verification); err != nil {
		return errs.InternalServerError("failed to create verification token")
	}

	// Send email async
	go func() {
		cfg := config.GetConfig()
		smtpCfg := utils.SMTPConfig{
			Host: cfg.SMTPHost, Port: cfg.SMTPPort,
			Email: cfg.SMTPEmail, Password: cfg.SMTPPassword,
			AppPort: cfg.PORT,
		}
		if err := utils.SendEmail(user.Email, token, smtpCfg); err != nil {
			log.Printf("Failed to send verification email to %s: %v", user.Email, err)
		}
	}()

	fmt.Printf("Verification token for %s: %s\n", user.Email, token)

	return nil
}
