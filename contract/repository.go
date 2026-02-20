package contract

import "mini_jira/models"

type Repository struct {
	User              UserRepository
	RefreshToken      RefreshTokenRepository
	EmailVerification EmailVerificationRepository
}

type UserRepository interface {
	GetUser(Id uint) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
	CreateUser(User *models.User) error
	UpdateUser(Id uint, User *models.User) error
	UpdateUserStatus(Id uint, status string) error
	DeleteUser(Id uint) error
}

type RefreshTokenRepository interface {
	Create(refreshToken *models.RefreshToken) error
	GetByTokenHash(tokenHash string) (*models.RefreshToken, error)
	RevokeByTokenHash(tokenHash string) error
	RevokeByUserId(userId uint) error
	RevokeByUserAndDevice(userId uint, deviceId string) error
}

type EmailVerificationRepository interface {
	Create(verification *models.EmailsVerification) error
	GetByToken(token string) (*models.EmailsVerification, error)
	Delete(id uint) error
	DeleteByUserID(userID uint) error
}
