package service

import (
	"log"
	"mini_jira/config"
	"mini_jira/contract"
	"mini_jira/dto"
	"mini_jira/models"
	errs "mini_jira/pkg/error"
	"mini_jira/pkg/token"
	"mini_jira/utils"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepository         contract.UserRepository
	TokenManager           *token.Manager
	RefreshTokenRepo       contract.RefreshTokenRepository
	EmailVerificationRepo  contract.EmailVerificationRepository
	RefreshTokenHashSecret string
}

func ImplUserService(repo contract.UserRepository, refreshRepo contract.RefreshTokenRepository, emailVerifRepo contract.EmailVerificationRepository, tokenManager *token.Manager, refreshHashSecret string) contract.UserService {
	return &UserService{
		UserRepository:         repo,
		TokenManager:           tokenManager,
		RefreshTokenRepo:       refreshRepo,
		EmailVerificationRepo:  emailVerifRepo,
		RefreshTokenHashSecret: refreshHashSecret,
	}
}

func (s *UserService) Login(payload *dto.UserRequest) (*dto.LoginResponse, error) {
	if !IsValidEmail(payload.Email) {
		return nil, errs.BadRequest("invalid email format")
	}

	if payload.Email == "" && payload.Password == "" {
		return nil, errs.BadRequest("email and password are required")
	}

	user, err := s.UserRepository.GetUserByEmail(payload.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.Unauthorized("invalid email or password")
		}
		return nil, errs.InternalServerError("failed to get user")
	}

	if user.Status != "active" {
		return nil, errs.Unauthorized("user is " + user.Status)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		return nil, errs.Unauthorized("invalid email or password")
	}

	if s.TokenManager == nil {
		return nil, errs.InternalServerError("token manager is not initialized")
	}
	if s.RefreshTokenHashSecret == "" {
		return nil, errs.InternalServerError("refresh token hash secret is not configured")
	}

	accessToken, err := s.TokenManager.NewAccessToken(user.Id, user.Username, string(user.Role))
	if err != nil {
		return nil, errs.InternalServerError("failed to create access token")
	}
	refreshToken, err := s.TokenManager.NewRefreshToken(user.Id)
	if err != nil {
		return nil, errs.InternalServerError("failed to create refresh token")
	}
	if s.RefreshTokenRepo == nil {
		return nil, errs.InternalServerError("refresh token repository is not initialized")
	}
	deviceId := strings.TrimSpace(payload.DeviceId)
	if deviceId == "" {
		deviceId = "default"
	}
	if err := s.RefreshTokenRepo.RevokeByUserAndDevice(user.Id, deviceId); err != nil {
		return nil, errs.InternalServerError("failed to revoke refresh tokens for device")
	}
	refreshExpiresAt := time.Now().UTC().Add(s.TokenManager.RefreshTTL())
	refreshTokenHash := utils.HMACSHA256Hex(refreshToken, s.RefreshTokenHashSecret)
	if err := s.RefreshTokenRepo.Create(&models.RefreshToken{
		UserID:    user.Id,
		TokenHash: refreshTokenHash,
		DeviceId:  deviceId,
		ExpiresAt: refreshExpiresAt,
	}); err != nil {
		return nil, errs.InternalServerError("failed to store refresh token")
	}

	response := &dto.LoginResponse{
		StatusCode: http.StatusOK,
		Message:    "Login successful",
		Data: dto.LoginData{
			Id:       user.Id,
			Name:     user.Name,
			Username: user.Username,
			Email:    user.Email,
			Role:     string(user.Role),
			Status:   user.Status,
		},
		Tokens: dto.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    int64(s.TokenManager.AccessTTL().Seconds()),
		},
	}

	return response, nil
}

func (s *UserService) RefreshToken(payload *dto.TokenRequest) (*dto.TokenResponse, error) {
	if payload == nil || payload.RefreshToken == "" {
		return nil, errs.BadRequest("refresh token is required")
	}
	if s.TokenManager == nil {
		return nil, errs.InternalServerError("token manager is not initialized")
	}
	if s.RefreshTokenRepo == nil {
		return nil, errs.InternalServerError("refresh token repository is not initialized")
	}
	if s.RefreshTokenHashSecret == "" {
		return nil, errs.InternalServerError("refresh token hash secret is not configured")
	}

	claims, err := s.TokenManager.VerifyRefreshToken(payload.RefreshToken)
	if err != nil {
		return nil, errs.Unauthorized("invalid or expired refresh token")
	}

	refreshTokenHash := utils.HMACSHA256Hex(payload.RefreshToken, s.RefreshTokenHashSecret)
	storedToken, err := s.RefreshTokenRepo.GetByTokenHash(refreshTokenHash)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.Unauthorized("invalid or expired refresh token")
		}
		return nil, errs.InternalServerError("failed to get refresh token")
	}
	if storedToken.RevokedAt != nil || storedToken.ExpiresAt.Before(time.Now().UTC()) {
		return nil, errs.Unauthorized("invalid or expired refresh token")
	}
	if storedToken.UserID != claims.UserID {
		return nil, errs.Unauthorized("invalid or expired refresh token")
	}
	deviceId := strings.TrimSpace(payload.DeviceId)
	if deviceId == "" {
		deviceId = "default"
	}
	if storedToken.DeviceId != deviceId {
		return nil, errs.Unauthorized("invalid or expired refresh token")
	}

	user, err := s.UserRepository.GetUser(claims.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.Unauthorized("invalid or expired refresh token")
		}
		return nil, errs.InternalServerError("failed to get user")
	}

	accessToken, err := s.TokenManager.NewAccessToken(user.Id, user.Username, string(user.Role))
	if err != nil {
		return nil, errs.InternalServerError("failed to create access token")
	}
	newRefreshToken, err := s.TokenManager.NewRefreshToken(user.Id)
	if err != nil {
		return nil, errs.InternalServerError("failed to create refresh token")
	}

	if err := s.RefreshTokenRepo.RevokeByTokenHash(refreshTokenHash); err != nil {
		return nil, errs.InternalServerError("failed to revoke refresh token")
	}
	newRefreshTokenHash := utils.HMACSHA256Hex(newRefreshToken, s.RefreshTokenHashSecret)
	if err := s.RefreshTokenRepo.Create(&models.RefreshToken{
		UserID:    claims.UserID,
		TokenHash: newRefreshTokenHash,
		DeviceId:  deviceId,
		ExpiresAt: time.Now().UTC().Add(s.TokenManager.RefreshTTL()),
	}); err != nil {
		return nil, errs.InternalServerError("failed to store refresh token")
	}

	return &dto.TokenResponse{
		StatusCode: http.StatusOK,
		Message:    "Token refreshed successfully",
		Tokens: dto.TokenPair{
			AccessToken:  accessToken,
			RefreshToken: newRefreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    int64(s.TokenManager.AccessTTL().Seconds()),
		},
	}, nil
}

func (s *UserService) Logout(payload *dto.TokenRequest) (*dto.BasicResponse, error) {
	if payload == nil || payload.RefreshToken == "" {
		return nil, errs.BadRequest("refresh token is required")
	}
	if s.TokenManager == nil {
		return nil, errs.InternalServerError("token manager is not initialized")
	}
	if s.RefreshTokenRepo == nil {
		return nil, errs.InternalServerError("refresh token repository is not initialized")
	}
	if s.RefreshTokenHashSecret == "" {
		return nil, errs.InternalServerError("refresh token hash secret is not configured")
	}

	if _, err := s.TokenManager.VerifyRefreshToken(payload.RefreshToken); err != nil {
		return nil, errs.Unauthorized("invalid or expired refresh token")
	}
	refreshTokenHash := utils.HMACSHA256Hex(payload.RefreshToken, s.RefreshTokenHashSecret)
	if err := s.RefreshTokenRepo.RevokeByTokenHash(refreshTokenHash); err != nil {
		return nil, errs.InternalServerError("failed to revoke refresh token")
	}

	return &dto.BasicResponse{
		StatusCode: http.StatusOK,
		Message:    "Logout successful",
	}, nil
}

func (s *UserService) Register(payload *dto.UserRequest) (*dto.UserResponse, error) {
	if !IsValidUsername(payload.Username) {
		return nil, errs.BadRequest("invalid username format")
	}
	if !IsValidEmail(payload.Email) {
		return nil, errs.BadRequest("invalid email format")
	}
	if !IsValidPassword(payload.Password) {
		return nil, errs.BadRequest("password need at least 8 characters with uppercase, lowercase, number, and special character")
	}

	usernameExists, err := s.UserRepository.GetUserByUsername(payload.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errs.InternalServerError("failed to check existing username")
	}
	if err == nil && usernameExists != nil {
		return nil, errs.Conflict("username already exists")
	}
	emailExists, err := s.UserRepository.GetUserByEmail(payload.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errs.InternalServerError("failed to check existing email")
	}
	if err == nil && emailExists != nil {
		return nil, errs.Conflict("email already exists")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.InternalServerError("failed to hash password")
	}

	role := models.RoleMember

	user := &models.User{
		Name:     payload.Name,
		Username: payload.Username,
		Email:    payload.Email,
		Password: string(hashPassword),
		Role:     role,
		Status: "Inactive",
	}

	err = s.UserRepository.CreateUser(user)
	if err != nil {
		return nil, errs.InternalServerError("failed to create user")
	}

	response := &dto.UserResponse{
		StatusCode: http.StatusCreated,
		Message:    "user registered successfully",
		Data: dto.UserData{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	token, err := utils.GenerateToken()
	if err != nil {
		return nil, errs.InternalServerError("failed to generate token")
	}

	verification := &models.EmailsVerification{
		UserID:    user.Id,
		Email:     user.Email,
		Token:     token,
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}
	if err := s.EmailVerificationRepo.Create(verification); err != nil {
		return nil, errs.InternalServerError("failed to create verification token")
	}

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

	return response, nil
}

func (s *UserService) Update(payload *dto.UserRequest) (*dto.UserUpdateResponse, error) {
	user, err := s.UserRepository.GetUser(payload.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NotFound("User not found")
		}
		return nil, errs.InternalServerError("Failed to get user")
	}

	if payload.Name != "" {
		user.Name = payload.Name
	}

	if payload.Username != "" {
		if !IsValidUsername(payload.Username) {
			return nil, errs.BadRequest("Invalid username format (3-30 chars, lowercase alphanumeric)")
		}
		user.Username = payload.Username
	}

	if payload.Role != "" {
		if payload.RequestingRole != string(models.RoleAdmin) {
			return nil, errs.Forbidden("Only admins can change user roles")
		}
		if payload.Role == string(models.RoleAdmin) {
			user.Role = models.RoleAdmin
		} else if payload.Role == string(models.RoleProjectManager) {
			user.Role = models.RoleProjectManager
		} else if payload.Role == string(models.RoleMember) {
			user.Role = models.RoleMember
		} else {
			return nil, errs.BadRequest("Invalid role")
		}
	}

	if payload.Email != "" {
		if !IsValidEmail(payload.Email) {
			return nil, errs.BadRequest("Invalid email format")
		}
		user.Email = payload.Email
	}

	if payload.AvatarUrl != "" {
		user.AvatarUrl = payload.AvatarUrl
	}

	if payload.NewPassword != "" {
		if payload.OldPassword == "" {
			return nil, errs.BadRequest("Old password is required")
		}
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.OldPassword))
		if err != nil {
			return nil, errs.Unauthorized("Old password is incorrect")
		}
		if !IsValidPassword(payload.NewPassword) {
			return nil, errs.BadRequest("Password must be at least 8 characters with uppercase, lowercase, number, and special character")
		}
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, errs.InternalServerError("Failed to hash password")
		}
		user.Password = string(hashPassword)
	}

	err = s.UserRepository.UpdateUser(payload.Id, user)
	if err != nil {
		return nil, errs.InternalServerError("Failed to update user")
	}

	response := &dto.UserUpdateResponse{
		StatusCode: http.StatusOK,
		Message:    "User updated successfully",
		Data: dto.UserUpdateData{
			Id:       user.Id,
			Username: user.Username,
			Status:   user.Status,
		},
	}

	return response, nil
}

func (s *UserService) Delete(userId uint) error {
	_, err := s.UserRepository.GetUser(userId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errs.NotFound("User not found")
		}
		return errs.InternalServerError("Failed to get user")
	}

	err = s.UserRepository.DeleteUser(userId)
	if err != nil {
		return errs.InternalServerError("Failed to delete user")
	}

	return nil
}

func (s *UserService) GetById(userId uint) (*dto.UserResponse, error) {
	user, err := s.UserRepository.GetUser(userId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NotFound("User not found")
		}
		return nil, errs.InternalServerError("Failed to get user by id")
	}

	response := &dto.UserResponse{
		StatusCode: http.StatusOK,
		Message:    "Successfully retrieved user data",
		Data: dto.UserData{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Role:      string(user.Role),
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}
	return response, nil
}

func (s *UserService) GetUserByUsername(username string) (*dto.UserResponse, error) {
	user, err := s.UserRepository.GetUserByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NotFound("User not found")
		}
		return nil, errs.InternalServerError("Failed to get user by username")
	}

	response := &dto.UserResponse{
		StatusCode: http.StatusOK,
		Message:    "Successfully retrieved user data",
		Data: dto.UserData{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Role:      string(user.Role),
			AvatarUrl: user.AvatarUrl,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}
	return response, nil
}

func (s *UserService) GetUserByEmail(email string) (*dto.UserResponse, error) {
	user, err := s.UserRepository.GetUserByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NotFound("User not found")
		}
		return nil, errs.InternalServerError("Failed to get user by email")
	}

	response := &dto.UserResponse{
		StatusCode: http.StatusOK,
		Message:    "Successfully retrieved user data",
		Data: dto.UserData{
			Id:        user.Id,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Role:      string(user.Role),
			AvatarUrl: user.AvatarUrl,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}
	return response, nil
}

func (s *UserService) GetAll() ([]*dto.UserAllResponse, error) {
	users, err := s.UserRepository.GetAllUsers()
	if err != nil {
		return nil, errs.InternalServerError("Failed to get all users")
	}

	var response []*dto.UserAllResponse
	for _, user := range users {
		response = append(response, &dto.UserAllResponse{
			Data: dto.UserData{
				Id:        user.Id,
				Name:      user.Name,
				Username:  user.Username,
				Email:     user.Email,
				Role:      string(user.Role),
				Status:    user.Status,
				AvatarUrl: user.AvatarUrl,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			},
		})
	}
	return response, nil
}

func (s *UserService) UpdateStatus(payload *dto.UserRequest) (*dto.UserUpdateResponse, error) {
	user, err := s.UserRepository.GetUser(payload.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NotFound("User not found")
		}
		return nil, errs.InternalServerError("Failed to get user")
	}

	if payload.Status != "active" && payload.Status != "inactive" && payload.Status != "suspended" {
		return nil, errs.BadRequest("Invalid status")
	}

	user.Status = payload.Status

	err = s.UserRepository.UpdateUser(payload.Id, user)
	if err != nil {
		return nil, errs.InternalServerError("Failed to update user status")
	}

	response := &dto.UserUpdateResponse{
		StatusCode: http.StatusOK,
		Message:    "User status updated successfully",
		Data: dto.UserUpdateData{
			Id:       user.Id,
			Username: user.Username,
		},
	}

	return response, nil
}

func (s *UserService) ResetPassword(payload *dto.UserRequest) (*dto.BasicResponse, error) {
	user, err := s.UserRepository.GetUser(payload.Id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NotFound("User not found")
		}
		return nil, errs.InternalServerError("Failed to get user by email")
	}

	if payload.NewPassword == "" {
		return nil, errs.BadRequest("New password is required")
	}

	if !IsValidPassword(payload.NewPassword) {
		return nil, errs.BadRequest("Invalid password format")
	}

	if payload.OldPassword == "" {
		return nil, errs.BadRequest("Old password is required")
	}

	if !IsValidPassword(payload.OldPassword) {
		return nil, errs.BadRequest("Invalid password format")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.OldPassword)); err != nil {
		return nil, errs.BadRequest("Old password is incorrect")
	}

	if payload.NewPassword == payload.OldPassword {
		return nil, errs.BadRequest("New password cannot be the same as old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.InternalServerError("Failed to hash password")
	}

	user.Password = string(hashedPassword)

	err = s.UserRepository.UpdateUser(user.Id, user)
	if err != nil {
		return nil, errs.InternalServerError("Failed to update user password")
	}

	response := &dto.BasicResponse{
		StatusCode: http.StatusOK,
		Message:    "Password updated successfully",
	}

	return response, nil
}

func IsValidUsername(username string) bool {
	if len(username) < 3 || len(username) > 30 {
		return false
	}
	regex := `^[a-z0-9]+([._][a-z0-9]+)*$`
	re, err := regexp.Compile(regex)
	if err != nil {
		return false
	}
	return re.MatchString(username)
}

func IsValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re, err := regexp.Compile(regex)
	if err != nil {
		return false
	}
	return re.MatchString(email)
}

func IsValidPassword(password string) bool {
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecialChar := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'"\\|,.<>\/?]`).MatchString(password)
	hasMinLen := len(password) >= 8

	return hasUpper && hasLower && hasNumber && hasSpecialChar && hasMinLen
}
