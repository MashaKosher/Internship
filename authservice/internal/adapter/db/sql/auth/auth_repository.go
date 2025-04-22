package auth

import (
	"authservice/internal/entity"

	db "authservice/pkg/client/sql"

	"gorm.io/gorm"
)

type AuthRepo struct {
	*gorm.DB
}

func New(db *gorm.DB) *AuthRepo {
	return &AuthRepo{db}
}

//

func (r *AuthRepo) CreateUser(user *entity.User) error {
	err := db.DB.Save(user).Error
	return err
}

func (r *AuthRepo) FindUserByName(username string) (entity.User, error) {
	var DBUser entity.User
	err := db.DB.Where("username = ?", username).First(&DBUser).Error
	if err != nil {
		return DBUser, err
	}
	return DBUser, nil
}

func (r *AuthRepo) FindUserById(userId int) (entity.User, error) {
	var DBUser entity.User
	if err := db.DB.Model(&DBUser).Where("id = ?", userId).First(&DBUser).Error; err != nil {
		return DBUser, err
	}
	return DBUser, nil
}

func (r *AuthRepo) ChangeUserPassword(userID int, newPassword string) error {
	result := db.DB.Model(&entity.User{}).Where("id = ?", userID).Update("password", newPassword)
	return result.Error
}

// /////////////////////////////////////////////

func CreateUser(user *entity.User) error {
	err := db.DB.Save(user).Error
	return err
}

func FindUserByName(username string) (entity.User, error) {
	var DBUser entity.User
	err := db.DB.Where("username = ?", username).First(&DBUser).Error
	if err != nil {
		return DBUser, err
	}
	return DBUser, nil
}

func FindUserById(userId int) (entity.User, error) {
	var DBUser entity.User
	if err := db.DB.Model(&DBUser).Where("id = ?", userId).First(&DBUser).Error; err != nil {
		return DBUser, err
	}
	return DBUser, nil
}

func ChangeUserPassword(userID int, newPassword string) error {
	result := db.DB.Model(&entity.User{}).Where("id = ?", userID).Update("password", newPassword)
	return result.Error
}
