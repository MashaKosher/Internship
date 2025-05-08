package auth

import (
	"authservice/internal/entity"
	"errors"

	"gorm.io/gorm"
)

type AuthRepo struct {
	*gorm.DB
}

func New(db *gorm.DB) *AuthRepo {
	return &AuthRepo{db}
}

const minUserID = 1

func (r *AuthRepo) CreateUser(user *entity.User) error {
	if user == nil {
		return entity.ErrUserIsNil
	}
	var count int64
	if err := r.DB.Model(&entity.User{}).Where("username = ?", user.Username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return entity.ErrUserAlreadyExists
	}
	return r.DB.Create(user).Error
}

func (r *AuthRepo) FindUserByName(username string) (entity.User, error) {
	if len(username) == 0 {
		return entity.User{}, entity.ErrUserNameCannotBeEmpty
	}
	var user entity.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, entity.ErrUserNotFoundInDB
		}
		return entity.User{}, err
	}
	return user, nil
}

func (r *AuthRepo) FindUserById(userID int) (entity.User, error) {
	if userID < minUserID {
		return entity.User{}, entity.ErrInvalidUserID
	}
	var user entity.User
	err := r.DB.First(&user, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, entity.ErrUserNotFoundInDB
	}
	return user, nil
}

func (r *AuthRepo) ChangeUserPassword(userID int, newPassword string) error {
	if userID < minUserID {
		return entity.ErrInvalidUserID
	}
	if newPassword == "" {
		return entity.ErrPasswordCannotBeEmpty
	}

	return r.DB.Transaction(func(tx *gorm.DB) error {
		var user entity.User
		if err := tx.First(&user, userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return entity.ErrUserNotFoundInDB
			}
			return err
		}

		return tx.Model(&user).Update("password", newPassword).Error
	})
}

func (r *AuthRepo) DeleteUser(userID int) error {
	if userID < minUserID {
		return entity.ErrInvalidUserID
	}

	result := r.DB.Delete(&entity.User{}, userID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entity.ErrUserNotFoundInDB
	}
	return nil
}
