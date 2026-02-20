package repository

import (
	"mini_jira/contract"
	"mini_jira/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func ImplUserRepository(db *gorm.DB) contract.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUser(Id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, Id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) CreateUser(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdateUser(Id uint, user *models.User) error {
	if err := r.db.Model(&models.User{}).Where("id = ?", Id).Updates(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdateUserStatus(Id uint, status string) error {
	return r.db.Model(&models.User{}).Where("id = ?", Id).Update("status", status).Error
}

func (r *UserRepository) DeleteUser(Id uint) error {
	if err := r.db.Delete(&models.User{}, Id).Error; err != nil {
		return err
	}
	return nil
}
