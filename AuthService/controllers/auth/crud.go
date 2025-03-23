package auth

import (
	db "authservice/config"
	models "authservice/models"
	"errors"

	"gorm.io/gorm"
)

func CreateUser(user *models.User) error {
	err := db.DB.Save(user).Error
	return err
}

func FindUser(user *models.User) (models.User, error) {
	var DBUser models.User
	err := db.DB.Where("username = ?", user.Username).First(&DBUser).Error
	if DBUser.ID == 0 {
		return DBUser, err
	}
	return DBUser, nil
}

func FindUserById(userId int) error {
	if err := db.DB.Model(&models.User{}).Where("id = ?", userId).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return nil
}
