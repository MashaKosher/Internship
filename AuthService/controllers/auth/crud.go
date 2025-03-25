package auth

import (
	db "authservice/config"
	models "authservice/models"
)

func CreateUser(user *models.User) error {
	err := db.DB.Save(user).Error
	return err
}

func FindUserByName(username string) (models.User, error) {
	var DBUser models.User
	err := db.DB.Where("username = ?", username).First(&DBUser).Error
	if err != nil {
		return DBUser, err
	}
	return DBUser, nil
}

func FindUserById(userId int) (models.User, error) {
	var DBUser models.User
	if err := db.DB.Model(&DBUser).Where("id = ?", userId).First(&DBUser).Error; err != nil {
		return DBUser, err
	}
	return DBUser, nil
}
