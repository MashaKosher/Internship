package service

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"authservice/internal/entity"
	"authservice/internal/logger"
	repo "authservice/internal/repository"

	"authservice/pkg/convert"
	"authservice/pkg/passwords"
	"authservice/pkg/tokens"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// @Summary		Sign up user
// @Description	Returns a message indicating the sign-up endpoint
// @Tags			SignUp
// @Accept			json
// @Produce		json
// @Param			models.User	body		models.User			true	"Sign up request body"
// @Success		200			{object}	models.UserResponse	"User successfully registered"
// @Router			/auth/sign-up/user [post]
func UserSignUp(c *fiber.Ctx) error {
	var user entity.User
	c.BodyParser(&user)
	logger.Logger.Info(fmt.Sprintf("User with UserName %s is trying to sign up", user.Username))

	// Veryfing income credentials
	if err := validate.Struct(user); err != nil {
		logger.Logger.Error("Body validation error: " + err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid data")
	}
	logger.Logger.Info("User credentials verified successfully")

	// Hashing Password to store in DB
	hashed, err := passwords.HashPassword(user.Password)
	if err != nil {
		logger.Logger.Error("Problems with hashing password: " + err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Problems with hashing password")
	}
	user.Password = hashed

	// Figure out if user exists
	if _, err := repo.FindUserByName(user.Username); err == nil {
		logger.Logger.Error("User with such username already exists: " + user.Username)
		return fiber.NewError(fiber.StatusBadRequest, "User with such username already exists")
	}

	// Adding User to DB
	if err := repo.CreateUser(&user); err != nil {
		logger.Logger.Error("Problem with creating User with UserName: " + user.Username)
		return fiber.NewError(fiber.StatusInternalServerError, "Problem with adding user to DB")
	}
	logger.Logger.Info("User created with ID: " + fmt.Sprint(user.ID))

	// Creating Access Token
	accessToken, err := tokens.CreateToken(tokens.ACCESS_TOKEN, &user)
	if err != nil {
		logger.Logger.Error("Problem with creating Access JWT Token" + err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Access JWT Token")
	}
	c.Cookie(&fiber.Cookie{
		Name:  tokens.ACCESS_TOKEN,
		Value: accessToken,
	})
	logger.Logger.Info("Access JWT created successfully")

	// Creating Refresh Token
	refreshToken, err := tokens.CreateToken(tokens.REFRESH_TOKEN, &user)
	if err != nil {
		logger.Logger.Error("Problem with creating Refresh JWT Token" + err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Refresh JWT Token")
	}

	c.Cookie(&fiber.Cookie{
		Name:  tokens.REFRESH_TOKEN,
		Value: refreshToken,
	})
	logger.Logger.Info("Refresh JWT created successfully")

	return c.JSON(convert.ConvertUserToResponse(accessToken, refreshToken, &user))

}

// @Summary		Sign up user
// @Description	Returns a message indicating the sign-up endpoint
// @Tags			SignUp
// @Accept			json
// @Produce		json
// @Param			models.User	body		models.User			true	"Sign up request body"
// @Success		200			{object}	models.UserResponse	"User successfully registered"
// @Router			/auth/sign-up/admin [post]
func AdminSignUp(c *fiber.Ctx) error {
	var user entity.User
	c.BodyParser(&user)
	logger.Logger.Info(fmt.Sprintf("User with UserName %s is trying to sign up", user.Username))

	// Veryfing income credentials
	if err := validate.Struct(user); err != nil {
		logger.Logger.Error("Body validation error: " + err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid data")
	}
	logger.Logger.Info("User credentials verified successfully")

	// Hashing Password to store in DB
	hashed, err := passwords.HashPassword(user.Password)
	if err != nil {
		logger.Logger.Error("Problems with hashing password: " + err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Problems with hashing password")
	}
	user.Password = hashed
	user.Role = "admin"

	// Figure out if user exists
	if _, err := repo.FindUserByName(user.Username); err == nil {
		logger.Logger.Error("User with such username already exists: " + user.Username)
		return fiber.NewError(fiber.StatusBadRequest, "User with such username already exists")
	}

	// Adding User to DB
	if err := repo.CreateUser(&user); err != nil {
		logger.Logger.Error("Problem with creating User with UserName: " + user.Username)
		return fiber.NewError(fiber.StatusInternalServerError, "Problem with adding user to DB")
	}
	logger.Logger.Info("User created with ID: " + fmt.Sprint(user.ID))

	// Creating Access Token
	accessToken, err := tokens.CreateToken(tokens.ACCESS_TOKEN, &user)
	if err != nil {
		logger.Logger.Error("Problem with creating Access JWT Token" + err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Access JWT Token")
	}
	c.Cookie(&fiber.Cookie{
		Name:  tokens.ACCESS_TOKEN,
		Value: accessToken,
	})
	logger.Logger.Info("Access JWT created successfully")

	// Creating Refresh Token
	refreshToken, err := tokens.CreateToken(tokens.REFRESH_TOKEN, &user)
	if err != nil {
		logger.Logger.Error("Problem with creating Refresh JWT Token" + err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Refresh JWT Token")
	}

	c.Cookie(&fiber.Cookie{
		Name:  tokens.REFRESH_TOKEN,
		Value: refreshToken,
	})
	logger.Logger.Info("Refresh JWT created successfully")

	return c.JSON(convert.ConvertUserToResponse(accessToken, refreshToken, &user))

}

// @Summary		User login
// @Description	Returns a message indicating the login endpoint
// @Tags			Login
// @Accept			json
// @Produce		json
// @Param			models.User	body		models.User				true	"Login request body"
// @Success		200			{object}	models.UserResponse		"User successfully logged"
// @Failure		400			{object}	models.ErrorResponse	"Invalid Username or Password"
// @Failure		500			{object}	models.ErrorResponse	"Internal Server Error"
// @Router			/auth/login [post]
func Login(c *fiber.Ctx) error {
	var user entity.User
	c.BodyParser(&user)
	logger.Logger.Info(fmt.Sprintf("User with UserName %s is trying to log in", user.Username))

	// Veryfing income credentials
	if err := validate.Struct(user); err != nil {
		logger.Logger.Error("Body validation error: " + err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid data")
	}
	logger.Logger.Info("User credentials verified successfully")

	// Search for this User in DB
	DBUser, err := repo.FindUserByName(user.Username)
	if err != nil {
		logger.Logger.Error("No user with such Username: " + user.Username)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Username")
	}
	logger.Logger.Info("User found, his ID: " + fmt.Sprint(DBUser.ID))

	// Comparing hashed password from DB to raw password from credentials
	if err := passwords.ValidatePassword(DBUser.Password, user.Password); err != nil {
		logger.Logger.Error("Invalid password: " + user.Password)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Password")
	}
	logger.Logger.Info("User with ID " + fmt.Sprint(DBUser.ID) + " has correct password")

	// Creating Access Token
	accessToken, err := tokens.CreateToken(tokens.ACCESS_TOKEN, &DBUser)
	if err != nil {
		logger.Logger.Error("Problem with creating Access JWT Token" + err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Access JWT Token")
	}
	c.Cookie(&fiber.Cookie{
		Name:  tokens.ACCESS_TOKEN,
		Value: accessToken,
	})
	logger.Logger.Info("Access JWT created successfully")

	// Creating Refresh Token
	refreshToken, err := tokens.CreateToken(tokens.REFRESH_TOKEN, &DBUser)
	if err != nil {
		logger.Logger.Error("Problem with creating Refresh JWT Token" + err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Refresh JWT Token")
	}

	c.Cookie(&fiber.Cookie{
		Name:  tokens.REFRESH_TOKEN,
		Value: refreshToken,
	})
	logger.Logger.Info("Refresh JWT created successfully")

	return c.JSON(convert.ConvertUserToResponse(accessToken, refreshToken, &DBUser))
}

// @Summary		Verifying access Token
// @Description	Verifying access, extract sub and returns Token status. Clears the Cookies, if there any error
// @Tags			Token
// @Produce		json
// @Success		200	{object}	models.UserResponse		"Access Token is Valid"
// @Failure		400	{object}	models.ErrorResponse	"Invalid access Token or No such User"
// @Failure		500	{object}	models.ErrorResponse	"Internal Server Error"
// @Router			/auth/check-token [get]
func CheckToken(c *fiber.Ctx) error {
	var accessToken string = c.Cookies(tokens.ACCESS_TOKEN)

	// Access Token Verifying
	validatedToken, err := tokens.TokenVerify(accessToken)

	// If access Token is invalid we clear the cookie and throw error
	if err != nil {
		c.ClearCookie()
		logger.Logger.Error("Inavlid " + tokens.ACCESS_TOKEN + " Token: " + err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	logger.Logger.Info(tokens.ACCESS_TOKEN + " Token is valid")

	// Get Token type from Token
	tokenType, err := tokens.GetTypeFromValidatedToken(validatedToken)
	if err != nil {
		logger.Logger.Error(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Verifying Token type (it must be access)
	if err := tokens.VerifyTokenType(tokens.ACCESS_TOKEN, tokenType); err != nil {
		logger.Logger.Error(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	logger.Logger.Info("Token type is correct")

	// Extracting User ID from valid access Token
	userId, err := tokens.GetIdFromValidatedToken(validatedToken)
	if err != nil {
		logger.Logger.Error(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	logger.Logger.Info("User ID from " + tokens.ACCESS_TOKEN + " Token is: " + fmt.Sprint(userId))

	// Search User by Id in DB
	DBUser, err := repo.FindUserById(int(userId))
	if err != nil {
		c.ClearCookie()
		logger.Logger.Error("No User with such id: " + fmt.Sprint(userId))
		return fiber.NewError(fiber.StatusBadRequest, "No such User")
	}
	logger.Logger.Info("User found successfully")

	return c.JSON(convert.ConvertUserToResponse(accessToken, "", &DBUser))
}

// @Summary		Verifying refresh Token and returning Access
// @Description	Verifying access, extract sub and returns Token status. Clears the Cookies, if there any error
// @Tags			Token
// @Produce		json
// @Success		200	{object} models.UserResponse		"Refresh Token is Valid"
// @Failure		400	{object}	models.ErrorResponse	"Invalid access Token or No such User"
// @Failure		500	{object}	models.ErrorResponse	"Internal Server Error"
// @Router			/auth/refresh [get]
func Refresh(c *fiber.Ctx) error {
	var refreshToken string = c.Cookies(tokens.REFRESH_TOKEN)

	// Refresh Token Verifying
	validatedToken, err := tokens.TokenVerify(refreshToken)

	// If refresh Token is invalid we clear the cookie and throw error
	if err != nil {
		c.ClearCookie()
		logger.Logger.Error("Inavlid " + tokens.REFRESH_TOKEN + " Token: " + err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	logger.Logger.Info(tokens.REFRESH_TOKEN + " Token is valid")

	// Get Token type from Token
	tokenType, err := tokens.GetTypeFromValidatedToken(validatedToken)
	if err != nil {
		logger.Logger.Error(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Verifying Token type (it must be refresh)
	if err := tokens.VerifyTokenType(tokens.REFRESH_TOKEN, tokenType); err != nil {
		logger.Logger.Error(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	logger.Logger.Info("Token type is correct")

	// Extracting User ID from valid refresh Token
	userId, err := tokens.GetIdFromValidatedToken(validatedToken)
	if err != nil {
		logger.Logger.Error(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	logger.Logger.Info("User ID from " + tokens.REFRESH_TOKEN + " Token is: " + fmt.Sprint(userId))

	// Search User by Id in DB
	DBUser, err := repo.FindUserById(int(userId))
	if err != nil {
		c.ClearCookie()
		logger.Logger.Error("No User with such id: " + fmt.Sprint(userId))
		return fiber.NewError(fiber.StatusBadRequest, "No such User")
	}
	logger.Logger.Info("User found successfully")

	// Creating Access Token
	accessToken, err := tokens.CreateToken(tokens.ACCESS_TOKEN, &DBUser)
	if err != nil {
		logger.Logger.Error("Problem with creating Access JWT Token" + err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Access JWT Token")
	}
	c.Cookie(&fiber.Cookie{
		Name:  tokens.ACCESS_TOKEN,
		Value: accessToken,
	})
	logger.Logger.Info("Access JWT created successfully")

	return c.JSON(convert.ConvertUserToResponse(accessToken, refreshToken, &DBUser))
}

// ChangePassword godoc
// @Summary Change user password
// @Description Updates authenticated user's password after validation
// @Tags Change Password
// @Accept json
// @Produce json
// @Param request body entity.Password true "New password details"
// @Success 200 {object} entity.Response "Password changed successfully"
// @Router /auth/change-password [post]
func ChangePassword(c *fiber.Ctx) error {

	var newPassword entity.Password
	c.BodyParser(&newPassword)
	// logger.Logger.Info(fmt.Sprintf("User with UserName %s is trying to log in", user.Username))

	// Veryfing income credentials
	if err := validate.Struct(newPassword); err != nil {
		logger.Logger.Error("Body validation error: " + err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid data")
	}
	logger.Logger.Info("User credentials verified successfully")

	var accessToken string = c.Cookies(tokens.ACCESS_TOKEN)
	// Access Token Verifying
	validatedToken, err := tokens.TokenVerify(accessToken)

	// If access Token is invalid we clear the cookie and throw error
	if err != nil {
		c.ClearCookie()
		logger.Logger.Error("Inavlid " + tokens.ACCESS_TOKEN + " Token: " + err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	logger.Logger.Info(tokens.ACCESS_TOKEN + " Token is valid")

	// Get Token type from Token
	tokenType, err := tokens.GetTypeFromValidatedToken(validatedToken)
	if err != nil {
		logger.Logger.Error(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Verifying Token type (it must be access)
	if err := tokens.VerifyTokenType(tokens.ACCESS_TOKEN, tokenType); err != nil {
		logger.Logger.Error(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	logger.Logger.Info("Token type is correct")

	// Extracting User ID from valid access Token
	userId, err := tokens.GetIdFromValidatedToken(validatedToken)
	if err != nil {
		logger.Logger.Error(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	logger.Logger.Info("User ID from " + tokens.ACCESS_TOKEN + " Token is: " + fmt.Sprint(userId))

	newHashedPassword, err := passwords.HashPassword(newPassword.NewPassword)
	if err != nil {
		logger.Logger.Error("Problems with hashing password: " + err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Problems with hashing password")
	}

	if err = repo.ChangeUserPassword(userId, newHashedPassword); err != nil {
		logger.Logger.Error("Problems with updating password: " + err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Problems with updating password")
	}
	logger.Logger.Info("Password updated successfully")
	return c.JSON(entity.Response{Message: "Password updated successfully"})
}
