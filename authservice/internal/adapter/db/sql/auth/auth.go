package auth

import (
	"authservice/internal/entity"

	"gorm.io/gorm"
)

type AuthRepo struct {
	*gorm.DB
}

func New(db *gorm.DB) *AuthRepo {
	return &AuthRepo{db}
}

func (r *AuthRepo) CreateUser(user *entity.User) error {
	err := r.DB.Save(user).Error
	return err
}

func (r *AuthRepo) FindUserByName(username string) (entity.User, error) {
	var DBUser entity.User
	err := r.DB.Where("username = ?", username).First(&DBUser).Error
	if err != nil {
		return DBUser, err
	}
	return DBUser, nil
}

func (r *AuthRepo) FindUserById(userId int) (entity.User, error) {
	var DBUser entity.User
	if err := r.DB.Model(&DBUser).Where("id = ?", userId).First(&DBUser).Error; err != nil {
		return DBUser, err
	}
	return DBUser, nil
}

func (r *AuthRepo) ChangeUserPassword(userID int, newPassword string) error {
	result := r.DB.Model(&entity.User{}).Where("id = ?", userID).Update("password", newPassword)
	return result.Error
}
