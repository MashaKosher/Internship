package di

import (
	"authservice/internal/entity"
)

type Services struct {
	Auth AuthService
}

type (
	AuthService interface {
		Login(user entity.UserInDTO) (entity.LoginOutDTO, error)
		CheckAccessToken(accessToken string) (entity.LoginOutDTO, error)
		CheckRefreshToken(refreshToken string) (entity.LoginOutDTO, error)
		UserSignUp(user entity.User) (entity.LoginOutDTO, error)
		AdminSignUp(user entity.User) (entity.LoginOutDTO, error)
		ChangePassword(newPassword entity.Password, accessToken string) (entity.LoginOutDTO, error)
	}
)
