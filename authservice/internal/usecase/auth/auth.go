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

func checkToken(token string, tokenTpe tokens.TokenType, r repo.AuthRepo) (entity.User, error) {
	// Access Token Verifying
	validatedToken, err := tokens.TokenVerify(token)

	// If access Token is invalid we clear the cookie and throw error
	if err != nil {
		logger.Logger.Error("Inavlid " + string(tokenTpe) + " Token: " + err.Error())
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	logger.Logger.Info(string(tokenTpe) + " Token is valid")

	// Get Token type from Token
	tokenType, err := tokens.GetTypeFromValidatedToken(validatedToken)
	if err != nil {
		logger.Logger.Error(err.Error())
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Verifying Token type (it must be access)
	if err := tokens.VerifyTokenType(string(tokenTpe), tokenType); err != nil {
		logger.Logger.Error(err.Error())
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	logger.Logger.Info("Token type is correct")

	// Extracting User ID from valid access Token
	userId, err := tokens.GetIdFromValidatedToken(validatedToken)
	if err != nil {
		logger.Logger.Error(err.Error())
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	logger.Logger.Info("User ID from " + string(tokenTpe) + " Token is: " + fmt.Sprint(userId))

	// Search User by Id in DB
	user, err := r.FindUserById(int(userId))
	if err != nil {
		logger.Logger.Error("No User with such id: " + fmt.Sprint(userId))
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, "No such User")
	}
	logger.Logger.Info("User found successfully")

	return user, nil
}

func (uc *UseCase) CheckAccessToken(accessToken string) (LoginOutDTO, error) {
	user, err := checkToken(accessToken, tokens.ACCESS_TOKEN, uc.repo)
	if err != nil {
		return LoginOutDTO{}, err
	}
	return convertUserToLoginOutDTO(user, accessToken, ""), nil
}

func (uc *UseCase) CheckRefreshToken(refreshToken string) (LoginOutDTO, error) {
	user, err := checkToken(refreshToken, tokens.REFRESH_TOKEN, uc.repo)
	if err != nil {
		return LoginOutDTO{}, err
	}

	// Creating Access Token
	accessToken, err := tokens.CreateToken(tokens.ACCESS_TOKEN, &user)
	if err != nil {
		logger.Logger.Error("Problem with creating Access JWT Token" + err.Error())
		return LoginOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Access JWT Token")
	}
	logger.Logger.Info("Access JWT created successfully")
	return convertUserToLoginOutDTO(user, accessToken, refreshToken), nil
}

func (uc *UseCase) UserSignUp(user entity.User) (LoginOutDTO, error) {
	outUser, err := signUp(user, entity.UserRole, uc.repo)
	if err != nil {
		return LoginOutDTO{}, err
	}

	return outUser, nil
}

func (uc *UseCase) AdminSignUp(user entity.User) (LoginOutDTO, error) {
	outUser, err := signUp(user, entity.AdminRole, uc.repo)
	if err != nil {
		return LoginOutDTO{}, err
	}

	return outUser, nil
}

func signUp(user entity.User, userRole entity.Role, repo repo.AuthRepo) (LoginOutDTO, error) {
	// Hashing Password to store in DB
	hashed, err := passwords.HashPassword(user.Password)
	if err != nil {
		logger.Logger.Error("Problems with hashing password: " + err.Error())
		return LoginOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problems with hashing password")
	}
	user.Password = hashed
	user.Role = userRole

	// Figure out if user exists
	if _, err := repo.FindUserByName(user.Username); err == nil {
		logger.Logger.Error("User with such username already exists: " + user.Username)
		return LoginOutDTO{}, fiber.NewError(fiber.StatusBadRequest, "User with such username already exists")
	}

	// Adding User to DB
	if err := repo.CreateUser(&user); err != nil {
		logger.Logger.Error("Problem with creating User with UserName: " + user.Username)
		return LoginOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problem with adding user to DB")
	}
	logger.Logger.Info("User created with ID: " + fmt.Sprint(user.ID))

	// Creating Access Token
	accessToken, err := tokens.CreateToken(tokens.ACCESS_TOKEN, &user)
	if err != nil {
		logger.Logger.Error("Problem with creating Access JWT Token" + err.Error())
		return LoginOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Access JWT Token")
	}
	logger.Logger.Info("Access JWT created successfully")

	// Creating Refresh Token
	refreshToken, err := tokens.CreateToken(tokens.REFRESH_TOKEN, &user)
	if err != nil {
		logger.Logger.Error("Problem with creating Refresh JWT Token" + err.Error())
		return LoginOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Refresh JWT Token")
	}
	logger.Logger.Info("Refresh JWT created successfully")

	return convertUserToLoginOutDTO(user, accessToken, refreshToken), nil
}

func (uc *UseCase) ChangePassword(newPassword entity.Password, accessToken string) (LoginOutDTO, error) {

	user, err := checkToken(accessToken, tokens.ACCESS_TOKEN, uc.repo)
	if err != nil {
		return LoginOutDTO{}, err
	}

	newHashedPassword, err := passwords.HashPassword(newPassword.NewPassword)
	if err != nil {
		logger.Logger.Error("Problems with hashing password: " + err.Error())
		return LoginOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problems with hashing password")
	}

	if err = uc.repo.ChangeUserPassword(int(user.ID), newHashedPassword); err != nil {
		logger.Logger.Error("Problems with updating password: " + err.Error())
		return LoginOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problems with updating password")
	}

	return convertUserToLoginOutDTO(user, accessToken, ""), nil
}
