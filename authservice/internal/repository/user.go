package repository

import (
	// db "authservice/config"
	// models "authservice/models"

	"authservice/internal/entity"

	"authservice/internal/db"
)

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
