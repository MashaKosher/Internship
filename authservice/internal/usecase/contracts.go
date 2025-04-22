package usecase

import (
	"authservice/internal/entity"
	"authservice/internal/usecase/auth"
)

type (
	Auth interface {
		Login(user auth.UserInDTO) (auth.LoginOutDTO, error)
		CheckAccessToken(accessToken string) (auth.LoginOutDTO, error)
		CheckRefreshToken(refreshToken string) (auth.LoginOutDTO, error)
		UserSignUp(user entity.User) (auth.LoginOutDTO, error)
		AdminSignUp(user entity.User) (auth.LoginOutDTO, error)
		ChangePassword(newPassword entity.Password, accessToken string) (auth.LoginOutDTO, error)
	}
)
