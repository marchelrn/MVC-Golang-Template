package repository

import (
	"mini_jira/contract"

	"gorm.io/gorm"
)

func New(db *gorm.DB) *contract.Repository {
	return &contract.Repository{
		User:              ImplUserRepository(db),
		RefreshToken:      ImplRefreshTokenRepository(db),
		EmailVerification: ImplEmailVerificationRepository(db),
	}
}
