package auth

import (
	repo "authservice/internal/adapter/db/sql"
	"authservice/internal/di"

	"authservice/internal/entity"
	"authservice/pkg/passwords"
	"authservice/pkg/tokens"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type UseCase struct {
	repo    repo.AuthRepo
	logger  di.LoggerType
	RSAKeys di.RSAKeys
}

func New(r repo.AuthRepo, logger di.LoggerType, RSAKeys di.RSAKeys) *UseCase {
	return &UseCase{
		repo:    r,
		logger:  logger,
		RSAKeys: RSAKeys,
	}
}

func (uc *UseCase) Login(user entity.UserInDTO) (entity.LoginOutDTO, error) {
	// Search for this User in DB
	DBUser, err := uc.repo.FindUserByName(user.Username)
	if err != nil {
		uc.logger.Error("No user with such Username: " + user.Username)
		return entity.LoginOutDTO{}, fiber.NewError(fiber.StatusBadRequest, "Invalid Username")
	}
	uc.logger.Info("User found, his ID: " + fmt.Sprint(DBUser.ID))

	// Comparing hashed password from DB to raw password from credentials
	if err := passwords.ValidatePassword(DBUser.Password, user.Password); err != nil {
		uc.logger.Error("Invalid password: " + user.Password)
		return entity.LoginOutDTO{}, fiber.NewError(fiber.StatusBadRequest, "Invalid Password")
	}
	uc.logger.Info("User with ID " + fmt.Sprint(DBUser.ID) + " has correct password")

	// Creating Access Token
	accessToken, err := tokens.CreateToken(di.ACCESS_TOKEN, &DBUser, uc.logger, uc.RSAKeys)
	if err != nil {
		uc.logger.Error("Problem with creating Access JWT Token" + err.Error())
		return entity.LoginOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Access JWT Token")
	}

	refreshToken, err := tokens.CreateToken(di.REFRESH_TOKEN, &DBUser, uc.logger, uc.RSAKeys)
	if err != nil {
		uc.logger.Error("Problem with creating Refresh JWT Token" + err.Error())
		return entity.LoginOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Refresh JWT Token")
	}

	return convertUserToLoginOutDTO(DBUser, accessToken, refreshToken), nil

}

func convertUserToLoginOutDTO(user entity.User, accessToken, refreshToken string) entity.LoginOutDTO {
	return entity.LoginOutDTO{UserID: int(user.ID), UserName: user.Username, UserRole: string(user.Role), AccessToken: accessToken, RefreshToken: refreshToken}
}

func checkToken(token string, tokenTpe di.TokenType, r repo.AuthRepo, logger di.LoggerType, RSAKeys di.RSAKeys) (entity.User, error) {
	// Access Token Verifying
	validatedToken, err := tokens.TokenVerify(token, logger, RSAKeys)

	// If access Token is invalid we clear the cookie and throw error
	if err != nil {
		logger.Error("Inavlid " + string(tokenTpe) + " Token: " + err.Error())
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	logger.Info(string(tokenTpe) + " Token is valid")

	// Get Token type from Token
	tokenType, err := tokens.GetTypeFromValidatedToken(validatedToken)
	if err != nil {
		logger.Error(err.Error())
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Verifying Token type (it must be access)
	if err := tokens.VerifyTokenType(string(tokenTpe), tokenType); err != nil {
		logger.Error(err.Error())
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	logger.Info("Token type is correct")

	// Extracting User ID from valid access Token
	userId, err := tokens.GetIdFromValidatedToken(validatedToken)
	if err != nil {
		logger.Error(err.Error())
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	logger.Info("User ID from " + string(tokenTpe) + " Token is: " + fmt.Sprint(userId))

	// Search User by Id in DB
	user, err := r.FindUserById(int(userId))
	if err != nil {
		logger.Error("No User with such id: " + fmt.Sprint(userId))
		return entity.User{}, fiber.NewError(fiber.StatusBadRequest, "No such User")
	}
	logger.Info("User found successfully")

	return user, nil
}

func (uc *UseCase) CheckAccessToken(accessToken string) (entity.LoginOutDTO, error) {
	user, err := checkToken(accessToken, di.ACCESS_TOKEN, uc.repo, uc.logger, uc.RSAKeys)
	if err != nil {
		return entity.LoginOutDTO{}, err
	}
	return convertUserToLoginOutDTO(user, accessToken, ""), nil
}

func (uc *UseCase) CheckRefreshToken(refreshToken string) (entity.LoginOutDTO, error) {
	user, err := checkToken(refreshToken, di.REFRESH_TOKEN, uc.repo, uc.logger, uc.RSAKeys)
	if err != nil {
		return entity.LoginOutDTO{}, err
	}

	// Creating Access Token
	accessToken, err := tokens.CreateToken(di.ACCESS_TOKEN, &user, uc.logger, uc.RSAKeys)
	if err != nil {
		uc.logger.Error("Problem with creating Access JWT Token" + err.Error())
		return entity.LoginOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Access JWT Token")
	}
	uc.logger.Info("Access JWT created successfully")
	return convertUserToLoginOutDTO(user, accessToken, refreshToken), nil
}

func (uc *UseCase) UserSignUp(user entity.User) (entity.LoginOutDTO, error) {
	outUser, err := signUp(user, entity.UserRole, uc.repo, uc.logger, uc.RSAKeys)
	if err != nil {
		return entity.LoginOutDTO{}, err
	}

	return outUser, nil
}

func (uc *UseCase) AdminSignUp(user entity.User) (entity.LoginOutDTO, error) {
	outUser, err := signUp(user, entity.AdminRole, uc.repo, uc.logger, uc.RSAKeys)
	if err != nil {
		return entity.LoginOutDTO{}, err
	}

	return outUser, nil
}

func signUp(user entity.User, userRole entity.Role, repo repo.AuthRepo, logger di.LoggerType, RSAKeys di.RSAKeys) (entity.LoginOutDTO, error) {
	// Hashing Password to store in DB
	hashed, err := passwords.HashPassword(user.Password)
	if err != nil {
		logger.Error("Problems with hashing password: " + err.Error())
		return entity.LoginOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problems with hashing password")
	}
	user.Password = hashed
	user.Role = userRole

	// Figure out if user exists
	if _, err := repo.FindUserByName(user.Username); err == nil {
		logger.Error("User with such username already exists: " + user.Username)
		return entity.LoginOutDTO{}, fiber.NewError(fiber.StatusBadRequest, "User with such username already exists")
	}

	// Adding User to DB
	if err := repo.CreateUser(&user); err != nil {
		logger.Error("Problem with creating User with UserName: " + user.Username)
		return entity.LoginOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problem with adding user to DB")
	}
	logger.Info("User created with ID: " + fmt.Sprint(user.ID))

	// Creating Access Token
	accessToken, err := tokens.CreateToken(di.ACCESS_TOKEN, &user, logger, RSAKeys)
	if err != nil {
		logger.Error("Problem with creating Access JWT Token" + err.Error())
		return entity.LoginOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Access JWT Token")
	}
	logger.Info("Access JWT created successfully")

	// Creating Refresh Token
	refreshToken, err := tokens.CreateToken(di.REFRESH_TOKEN, &user, logger, RSAKeys)
	if err != nil {
		logger.Error("Problem with creating Refresh JWT Token" + err.Error())
		return entity.LoginOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Refresh JWT Token")
	}
	logger.Info("Refresh JWT created successfully")

	return convertUserToLoginOutDTO(user, accessToken, refreshToken), nil
}

func (uc *UseCase) ChangePassword(newPassword entity.Password, accessToken string) (entity.LoginOutDTO, error) {

	user, err := checkToken(accessToken, di.ACCESS_TOKEN, uc.repo, uc.logger, uc.RSAKeys)
	if err != nil {
		return entity.LoginOutDTO{}, err
	}

	newHashedPassword, err := passwords.HashPassword(newPassword.NewPassword)
	if err != nil {
		uc.logger.Error("Problems with hashing password: " + err.Error())
		return entity.LoginOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problems with hashing password")
	}

	if err = uc.repo.ChangeUserPassword(int(user.ID), newHashedPassword); err != nil {
		uc.logger.Error("Problems with updating password: " + err.Error())
		return entity.LoginOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problems with updating password")
	}

	return convertUserToLoginOutDTO(user, accessToken, ""), nil
}
