package di

import (
	"authservice/internal/entity"
)

type Services struct {
	Auth AuthService
}

type (
	AuthService interface {
		Login(user entity.UserInDTO) (entity.UserOutDTO, error)
		CheckAccessToken(accessToken string) (entity.UserOutDTO, error)
		CheckRefreshToken(refreshToken string) (entity.UserOutDTO, error)
		UserSignUp(user entity.User) (entity.UserOutDTO, error)
		AdminSignUp(user entity.User) (entity.UserOutDTO, error)
		ChangePassword(newPassword entity.Password, accessToken string) (entity.UserOutDTO, error)
	}
)
