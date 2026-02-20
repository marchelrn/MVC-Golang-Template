package repository

import (
	"mini_jira/contract"
	"mini_jira/models"
	"time"

	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func ImplRefreshTokenRepository(db *gorm.DB) contract.RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(refreshToken *models.RefreshToken) error {
	return r.db.Create(refreshToken).Error
}

func (r *RefreshTokenRepository) GetByTokenHash(tokenHash string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	if err := r.db.Where("token_hash = ?", tokenHash).First(&refreshToken).Error; err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *RefreshTokenRepository) RevokeByTokenHash(tokenHash string) error {
	now := time.Now().UTC()
	return r.db.Model(&models.RefreshToken{}).
		Where("token_hash = ? AND revoked_at IS NULL", tokenHash).
		Update("revoked_at", now).Error
}

func (r *RefreshTokenRepository) RevokeByUserId(userId uint) error {
	now := time.Now().UTC()
	return r.db.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userId).
		Update("revoked_at", now).Error
}

func (r *RefreshTokenRepository) RevokeByUserAndDevice(userId uint, deviceId string) error {
	now := time.Now().UTC()
	return r.db.Model(&models.RefreshToken{}).
		Where("user_id = ? AND device_id = ? AND revoked_at IS NULL", userId, deviceId).
		Update("revoked_at", now).Error
}
