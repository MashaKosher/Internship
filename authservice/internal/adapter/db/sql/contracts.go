package sql

import "authservice/internal/entity"

type (
	// TranslationRepo -.
	AuthRepo interface {
		CreateUser(user *entity.User) error
		FindUserByName(username string) (entity.User, error)
		FindUserById(userId int) (entity.User, error)
		ChangeUserPassword(userID int, newPassword string) error
	}
)
