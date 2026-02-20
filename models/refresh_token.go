package models

import "time"

type RefreshToken struct {
	Id        uint       `gorm:"primaryKey;autoIncrement"`
	UserID    uint       `gorm:"not null;index"`
	TokenHash string     `gorm:"type:text;uniqueIndex;not null"`
	DeviceId  string     `gorm:"type:varchar(128);not null;index"`
	ExpiresAt time.Time  `gorm:"not null;index"`
	RevokedAt *time.Time `gorm:"index"`
	CreatedAt time.Time
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
