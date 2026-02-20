package repository

import (
	"mini_jira/contract"
	"mini_jira/models"

	"gorm.io/gorm"
)

type EmailVerificationRepository struct {
	db *gorm.DB
}

func ImplEmailVerificationRepository(db *gorm.DB) contract.EmailVerificationRepository {
	return &EmailVerificationRepository{db: db}
}

func (r *EmailVerificationRepository) Create(verification *models.EmailsVerification) error {
	return r.db.Create(verification).Error
}

func (r *EmailVerificationRepository) GetByToken(token string) (*models.EmailsVerification, error) {
	var verification models.EmailsVerification
	if err := r.db.Where("token = ?", token).First(&verification).Error; err != nil {
		return nil, err
	}
	return &verification, nil
}

func (r *EmailVerificationRepository) Delete(id uint) error {
	return r.db.Delete(&models.EmailsVerification{}, id).Error
}

func (r *EmailVerificationRepository) DeleteByUserID(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.EmailsVerification{}).Error
}
