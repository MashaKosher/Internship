package auth

import (
	repo "authservice/internal/adapter/db/sql"
	"authservice/internal/entity"
	"authservice/pkg/logger"
	"authservice/pkg/passwords"
	"authservice/pkg/tokens"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type UseCase struct {
	repo repo.AuthRepo
}

func New(r repo.AuthRepo) *UseCase {
	return &UseCase{
		repo: r,
	}
}

func (uc *UseCase) Login(user UserInDTO) (LoginOutDTO, error) {
	// Search for this User in DB
	DBUser, err := uc.repo.FindUserByName(user.Username)
	if err != nil {
		logger.Logger.Error("No user with such Username: " + user.Username)
		return LoginOutDTO{}, fiber.NewError(fiber.StatusBadRequest, "Invalid Username")
	}
	logger.Logger.Info("User found, his ID: " + fmt.Sprint(DBUser.ID))

	// Comparing hashed password from DB to raw password from credentials
	if err := passwords.ValidatePassword(DBUser.Password, user.Password); err != nil {
		logger.Logger.Error("Invalid password: " + user.Password)
		return LoginOutDTO{}, fiber.NewError(fiber.StatusBadRequest, "Invalid Password")
	}
	logger.Logger.Info("User with ID " + fmt.Sprint(DBUser.ID) + " has correct password")

	// Creating Access Token
	accessToken, err := tokens.CreateToken(tokens.ACCESS_TOKEN, &DBUser)
	if err != nil {
		logger.Logger.Error("Problem with creating Access JWT Token" + err.Error())
		return LoginOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Access JWT Token")
	}

	refreshToken, err := tokens.CreateToken(tokens.REFRESH_TOKEN, &DBUser)
	if err != nil {
		logger.Logger.Error("Problem with creating Refresh JWT Token" + err.Error())
		return LoginOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Refresh JWT Token")
	}

	return convertUserToLoginOutDTO(DBUser, accessToken, refreshToken), nil

}

func convertUserToLoginOutDTO(user entity.User, accessToken, refreshToken string) LoginOutDTO {
	return LoginOutDTO{UserID: int(user.ID), UserName: user.Username, UserRole: string(user.Role), AccessToken: accessToken, RefreshToken: refreshToken}
}
