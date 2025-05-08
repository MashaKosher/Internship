package auth

import (
	repo "authservice/internal/adapter/db/sql"
	"authservice/internal/di"

	"authservice/internal/entity"
	"authservice/pkg/passwords"
	"authservice/pkg/tokens"
	"fmt"

	kafkaRepo "authservice/internal/adapter/kafka"

	"github.com/gofiber/fiber/v2"
)

type UseCase struct {
	repo    repo.AuthRepo
	logger  di.LoggerType
	RSAKeys di.RSAKeys
	Bus     kafkaRepo.SignUpProducer
	Cache   di.Cache
}

func New(r repo.AuthRepo, logger di.LoggerType, RSAKeys di.RSAKeys, signUpProducer kafkaRepo.SignUpProducer, cache di.Cache) *UseCase {
	return &UseCase{
		repo:    r,
		logger:  logger,
		RSAKeys: RSAKeys,
		Bus:     signUpProducer,
		Cache:   cache,
	}
}

func (uc *UseCase) Login(user entity.UserInDTO) (entity.UserOutDTO, error) {
	// Search for this User in DB
	DBUser, err := uc.repo.FindUserByName(user.Username)
	if err != nil {
		if err == entity.ErrUserNameCannotBeEmpty {
			uc.logger.Error("Invalid data to find user by name in DB: " + user.Username)
			return entity.UserOutDTO{}, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if err == entity.ErrUserNotFoundInDB {
			uc.logger.Error(err.Error())
			return entity.UserOutDTO{}, fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		uc.logger.Error("Some DB Error while finding user by name: " + user.Username)
		return entity.UserOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Some problems while finding user in DB")
	}
	uc.logger.Info("User found, his ID: " + fmt.Sprint(DBUser.ID))

	// Comparing hashed password from DB to raw password from credentials
	if err := passwords.ValidatePassword(DBUser.Password, user.Password); err != nil {
		uc.logger.Error("Invalid password: " + user.Password)
		return entity.UserOutDTO{}, fiber.NewError(fiber.StatusBadRequest, "Invalid Password")
	}
	uc.logger.Info("User with ID " + fmt.Sprint(DBUser.ID) + " has correct password")

	// Creating Access Token
	accessToken, err := tokens.CreateToken(di.ACCESS_TOKEN, &DBUser, uc.logger, uc.RSAKeys)
	if err != nil {
		uc.logger.Error("Problem with creating Access JWT Token" + err.Error())
		return entity.UserOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Access JWT Token")
	}
	// Adding access token to cache
	uc.Cache.Token.AddTokenToCache(accessToken, int(DBUser.ID))

	// Creating Access Token
	refreshToken, err := tokens.CreateToken(di.REFRESH_TOKEN, &DBUser, uc.logger, uc.RSAKeys)
	if err != nil {
		uc.logger.Error("Problem with creating Refresh JWT Token" + err.Error())
		return entity.UserOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Refresh JWT Token")
	}
	return DBUser.ToDTO(accessToken, refreshToken), nil

}

func checkToken(token string, mustTokenType di.TokenType, r repo.AuthRepo, logger di.LoggerType, RSAKeys di.RSAKeys, cache di.Cache) (entity.User, error) {

	mustValidate := true
	var err error
	var userId int

	if mustTokenType == di.ACCESS_TOKEN {
		userId, err = cache.Token.GetTokenFromCache(token)
		if err == nil {
			mustValidate = false
			logger.Info("Token found in cache")
		}
	}

	if mustValidate {
		logger.Info("Token not found in cache")
		// Access Token Verifying
		validatedToken, err := tokens.TokenVerify(token, logger, RSAKeys)
		// If Token is invalid we clear the cookie and throw error
		if err != nil {
			if err == entity.ErrTokenExpired {
				logger.Error("Token expired")
				return entity.User{}, fiber.NewError(fiber.StatusUnauthorized, err.Error())
			}
			if err == entity.ErrInvalidSigningMethod {
				logger.Error("Invalid signing methood")
				return entity.User{}, fiber.NewError(fiber.StatusForbidden, err.Error())
			}
			logger.Error("Some server error while parsing token: " + err.Error())
			return entity.User{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		logger.Info(string(mustTokenType) + " Token is valid")

		// Get Token type from Token
		tokenType, err := tokens.GetTypeFromValidatedToken(validatedToken)
		if err != nil {
			logger.Error(err.Error())
			return entity.User{}, fiber.NewError(fiber.StatusForbidden, err.Error())
		}

		// Verifying Token type (it must be same as mustTokenType in params)
		if err := tokens.VerifyTokenType(string(mustTokenType), tokenType); err != nil {
			logger.Error(err.Error())
			return entity.User{}, fiber.NewError(fiber.StatusForbidden, err.Error())
		}
		logger.Info("Token type is correct")

		// Extracting User ID from valid Token
		userId, err = tokens.GetIdFromValidatedToken(validatedToken)
		if err != nil {
			logger.Error(err.Error())
			return entity.User{}, fiber.NewError(fiber.StatusForbidden, err.Error())
		}
		logger.Info("User ID from " + string(mustTokenType) + " Token is: " + fmt.Sprint(userId))
	}

	// Search User by ID in DB
	user, err := r.FindUserById(int(userId))
	if err != nil {
		if err == entity.ErrInvalidUserID {
			logger.Error("Invalid data to update password in DB: " + err.Error())
			return entity.User{}, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if err == entity.ErrUserNotFoundInDB {
			logger.Error(err.Error())
			return entity.User{}, fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		logger.Error("Some DB Error while finding user by ID in DB: " + err.Error())
		return entity.User{}, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	logger.Info("User found successfully")

	return user, nil
}

func (uc *UseCase) CheckAccessToken(accessToken string) (entity.UserOutDTO, error) {
	user, err := checkToken(accessToken, di.ACCESS_TOKEN, uc.repo, uc.logger, uc.RSAKeys, uc.Cache)
	if err != nil {
		return entity.UserOutDTO{}, err
	}
	return user.ToDTO(accessToken, ""), nil
}

func (uc *UseCase) CheckRefreshToken(refreshToken string) (entity.UserOutDTO, error) {
	user, err := checkToken(refreshToken, di.REFRESH_TOKEN, uc.repo, uc.logger, uc.RSAKeys, uc.Cache)
	if err != nil {
		return entity.UserOutDTO{}, err
	}
	return user.ToDTO("", refreshToken), nil
}

func (uc *UseCase) CheckTokens(accessToken, refreshToken string) (entity.UserOutDTO, error) {
	user, err := checkToken(accessToken, di.ACCESS_TOKEN, uc.repo, uc.logger, uc.RSAKeys, uc.Cache)
	if err == nil {
		return user.ToDTO(accessToken, refreshToken), err
	}

	if err.Error() != entity.ErrTokenExpired.Error() {
		uc.logger.Error("Problem with creating Access JWT Token" + err.Error())
		return entity.UserOutDTO{}, err
	}

	user, err = checkToken(refreshToken, di.REFRESH_TOKEN, uc.repo, uc.logger, uc.RSAKeys, uc.Cache)
	if err != nil {
		if err == entity.ErrTokenExpired {
			return entity.UserOutDTO{}, fiber.NewError(fiber.StatusUnauthorized, "Token Expirtred")
		}
		return entity.UserOutDTO{}, err
	}

	// Creating Access Token
	accessToken, err = tokens.CreateToken(di.ACCESS_TOKEN, &user, uc.logger, uc.RSAKeys)
	if err != nil {
		uc.logger.Error("Problem with creating Access JWT Token" + err.Error())
		return entity.UserOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Access JWT Token")
	}
	uc.logger.Info("Access JWT created successfully")
	return user.ToDTO(accessToken, refreshToken), nil

}

func (uc *UseCase) UserSignUp(user entity.User, referalID int) (entity.UserOutDTO, error) {
	outUser, err := signUp(user, entity.UserRole, uc.repo, uc.logger, uc.RSAKeys, uc.Bus, referalID)
	if err != nil {
		return entity.UserOutDTO{}, err
	}
	return outUser, nil
}

func (uc *UseCase) AdminSignUp(user entity.User, referalID int) (entity.UserOutDTO, error) {
	outUser, err := signUp(user, entity.AdminRole, uc.repo, uc.logger, uc.RSAKeys, uc.Bus, referalID)
	if err != nil {
		return entity.UserOutDTO{}, err
	}
	return outUser, nil
}

func signUp(user entity.User, userRole entity.Role, repo repo.AuthRepo, logger di.LoggerType, RSAKeys di.RSAKeys, bus kafkaRepo.SignUpProducer, referalID int) (entity.UserOutDTO, error) {
	// Hashing Password to store in DB
	hashed, err := passwords.HashPassword(user.Password)
	if err != nil {
		logger.Error("Problems with hashing password: " + err.Error())
		return entity.UserOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problems with hashing password")
	}
	user.Password = hashed
	user.Role = userRole

	// Adding User to DB
	if err := repo.CreateUser(&user); err != nil {
		if err == entity.ErrUserIsNil {
			logger.Error("Invalid data to create user in DB: " + user.Username)
			return entity.UserOutDTO{}, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if err == entity.ErrUserAlreadyExists {
			logger.Error("User with such name already exists: " + user.Username)
			return entity.UserOutDTO{}, fiber.NewError(fiber.StatusConflict, err.Error())
		}

		logger.Error("Some DB Error while creating user: " + user.Username)
		return entity.UserOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problem with adding user to DB")
	}
	logger.Info("User created with ID: " + fmt.Sprint(user.ID))

	// Creating Access Token
	accessToken, err := tokens.CreateToken(di.ACCESS_TOKEN, &user, logger, RSAKeys)
	if err != nil {
		logger.Error("Problem with creating Access JWT Token" + err.Error())
		return entity.UserOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Access JWT Token")
	}
	logger.Info("Access JWT created successfully")

	// Creating Refresh Token
	refreshToken, err := tokens.CreateToken(di.REFRESH_TOKEN, &user, logger, RSAKeys)
	if err != nil {
		logger.Error("Problem with creating Refresh JWT Token" + err.Error())
		return entity.UserOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Refresh JWT Token")
	}
	logger.Info("Refresh JWT created successfully")

	// Send info to core service
	go bus.SendUserSignUpInfo(user.ToUserSignUpOutDTO(referalID))

	return user.ToDTO(accessToken, refreshToken), nil
}

func (uc *UseCase) ChangePassword(newPassword entity.Password, accessToken string) (entity.UserOutDTO, error) {

	user, err := checkToken(accessToken, di.ACCESS_TOKEN, uc.repo, uc.logger, uc.RSAKeys, uc.Cache)
	if err != nil {
		return entity.UserOutDTO{}, err
	}

	newHashedPassword, err := passwords.HashPassword(newPassword.NewPassword)
	if err != nil {
		uc.logger.Error("Problems with hashing password: " + err.Error())
		return entity.UserOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problems with hashing password")
	}

	if err = uc.repo.ChangeUserPassword(int(user.ID), newHashedPassword); err != nil {
		if err == entity.ErrInvalidUserID || err == entity.ErrPasswordCannotBeEmpty {
			uc.logger.Error("Invalid data to update user password in DB: " + err.Error())
			return entity.UserOutDTO{}, fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if err == entity.ErrUserNotFoundInDB {
			uc.logger.Error(err.Error())
			return entity.UserOutDTO{}, fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		uc.logger.Error("Some DB Error while updating user password: " + err.Error())
		return entity.UserOutDTO{}, fiber.NewError(fiber.StatusInternalServerError, "Problems with updating password")
	}

	return user.ToDTO(accessToken, ""), nil
}

func (uc *UseCase) DeleteUser(userId int) error {
	err := uc.repo.DeleteUser(userId)
	if err != nil {
		if err == entity.ErrInvalidUserID {
			uc.logger.Error("Invalid data to delet user from DB: " + err.Error())
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if err == entity.ErrUserNotFoundInDB {
			uc.logger.Error(err.Error())
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}

		uc.logger.Error("Some DB Error while deleting user: " + err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return nil
}
