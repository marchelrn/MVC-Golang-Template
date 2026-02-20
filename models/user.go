package models

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin          Role = "admin"
	RoleProjectManager Role = "project_manager"
	RoleMember         Role = "member"
)

type User struct {
	Id        uint   `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"type:varchar(100);not null"`
	Username  string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Password  string `gorm:"te:varchar(255);not null"`
	Email     string `gorm:"typype:varchar(100);uniqueIndex;not null"`
	Role      Role   `gorm:"type:varchar(50);not null;default:'member'"`
	AvatarUrl string `gorm:"type:varchar(255)"`
	Status    string `gorm:"type:varchar(50);not null;default:'active'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:""`
}

func (User) TableName() string {
	return "users"
}
