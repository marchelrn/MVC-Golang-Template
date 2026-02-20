package contract

import (
	"mini_jira/dto"
	"mini_jira/pkg/token"
)

type Service struct {
	User              UserService
	EmailVerification EmailVerificationService
	Token             *token.Manager
}

type UserService interface {
	Login(payload *dto.UserRequest) (*dto.LoginResponse, error)
	RefreshToken(payload *dto.TokenRequest) (*dto.TokenResponse, error)
	Logout(payload *dto.TokenRequest) (*dto.BasicResponse, error)
	Register(payload *dto.UserRequest) (*dto.UserResponse, error)
	Update(payload *dto.UserRequest) (*dto.UserUpdateResponse, error)
	Delete(userId uint) error
	GetById(userId uint) (*dto.UserResponse, error)
	GetUserByUsername(username string) (*dto.UserResponse, error)
	GetUserByEmail(email string) (*dto.UserResponse, error)
	GetAll() ([]*dto.UserAllResponse, error)
	UpdateStatus(payload *dto.UserRequest) (*dto.UserUpdateResponse, error)
	ResetPassword(payload *dto.UserRequest) (*dto.BasicResponse, error)
}

type EmailVerificationService interface {
	VerifyEmail(token string) error
	ResendVerification(email string) error
}
