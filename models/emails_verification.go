package models

import (
	"time"
)

type EmailsVerification struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"not null"`
	Email     string    `gorm:"type:varchar(255);not null"`
	Token     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	ExpiredAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (EmailsVerification) TableName() string {
	return "emails_verification"
}
