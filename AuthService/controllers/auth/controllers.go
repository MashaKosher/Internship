package auth

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	models "authservice/models"

	config "authservice/config"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// SignUp
// @Summary Sign up user
// @Description Returns a message indicating the sign-up endpoint
// @Tags authentication
// @Accept json
// @Produce json
// @Param models.User body models.User true "Sign up request body"
// @Success 200 {object} models.UserResponse "User successfully registered"
// @Router /auth/signup [post]
func SignUp(c *fiber.Ctx) error {
	var user models.User
	c.BodyParser(&user)
	config.Logger.Info(fmt.Sprintf("User with UserName %s is trying to sign up", user.Username))

	// Veryfing income credentials
	if err := validate.Struct(user); err != nil {
		config.Logger.Error("Body validation error: " + err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid data")
	}
	config.Logger.Info("User credentials verified successfully")

	// Hashing Password to store in DB
	hashed, err := HashPassword(user.Password)
	if err != nil {
		config.Logger.Error("Problems with hashing password: " + err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Problems with hashing password")
	}
	user.Password = hashed

	// Adding User to DB
	if err := CreateUser(&user); err != nil {
		config.Logger.Error("Problem with creating User with UserName: " + user.Username)
		return fiber.NewError(fiber.StatusInternalServerError, "Problem with adding user to DB")
	}
	config.Logger.Info("User created with ID: " + fmt.Sprint(user.ID))

	// Creating Access Token
	accessToken, err := CreateToken(ACCESS_TOKEN, &user)
	if err != nil {
		config.Logger.Error("Problem with creating Access JWT Token" + err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Access JWT Token")
	}
	c.Cookie(&fiber.Cookie{
		Name:  ACCESS_TOKEN,
		Value: accessToken,
	})
	config.Logger.Info("Access JWT created successfully")

	// Creating Refresh Token
	refreshToken, err := CreateToken(REFRESH_TOKEN, &user)
	if err != nil {
		config.Logger.Error("Problem with creating Refresh JWT Token" + err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Refresh JWT Token")
	}

	c.Cookie(&fiber.Cookie{
		Name:  REFRESH_TOKEN,
		Value: refreshToken,
	})
	config.Logger.Info("Refresh JWT created successfully")

	return c.JSON(ConvertUserToResponse(accessToken, refreshToken, &user))

}

// Login
// @Summary User login
// @Description Returns a message indicating the login endpoint
// @Tags authentication
// @Accept json
// @Produce json
// @Param models.User body models.User true "Login request body"
// @Success 200 {object} models.UserResponse "User successfully logged"
// @Failure 400 {object} models.ErrorResponse "Invalid Username or Password"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
	var user models.User
	c.BodyParser(&user)
	config.Logger.Info(fmt.Sprintf("User with UserName %s is trying to log in", user.Username))

	// Veryfing income credentials
	if err := validate.Struct(user); err != nil {
		config.Logger.Error("Body validation error: " + err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid data")
	}
	config.Logger.Info("User credentials verified successfully")

	// Search for this User in DB
	DBUser, err := FindUserByName(user.Username)
	if err != nil {
		config.Logger.Error("No user with such Username: " + user.Username)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Username")
	}
	config.Logger.Info("User found, his ID: " + fmt.Sprint(DBUser.ID))

	// Comparing hashed password from DB to raw password from credentials
	if err := ValidatePassword(DBUser.Password, user.Password); err != nil {
		config.Logger.Error("Invalid password: " + user.Password)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Password")
	}
	config.Logger.Info("User with ID " + fmt.Sprint(DBUser.ID) + " has correct password")

	// Creating Access Token
	accessToken, err := CreateToken(ACCESS_TOKEN, &DBUser)
	if err != nil {
		config.Logger.Error("Problem with creating Access JWT Token" + err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Access JWT Token")
	}
	c.Cookie(&fiber.Cookie{
		Name:  ACCESS_TOKEN,
		Value: accessToken,
	})
	config.Logger.Info("Access JWT created successfully")

	// Creating Refresh Token
	refreshToken, err := CreateToken(REFRESH_TOKEN, &DBUser)
	if err != nil {
		config.Logger.Error("Problem with creating Refresh JWT Token" + err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Refresh JWT Token")
	}

	c.Cookie(&fiber.Cookie{
		Name:  REFRESH_TOKEN,
		Value: refreshToken,
	})
	config.Logger.Info("Refresh JWT created successfully")

	return c.JSON(ConvertUserToResponse(accessToken, refreshToken, &DBUser))
}

// CheckToken check access Token From cookies
// @Summary Verifying access Token
// @Description Verifying access, extract sub and returns Token status. Clears the Cookies, if there any error
// @Tags authentication
// @Produce json
// @Success 200 {object} models.UserResponse "Access Token is Valid"
// @Failure 400 {object} models.ErrorResponse "Invalid access Token or No such User"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /auth/check-token [get]
func CheckToken(c *fiber.Ctx) error {
	var accessToken string = c.Cookies(ACCESS_TOKEN)

	// Access Token Verifying
	validatedToken, err := TokenVerify(accessToken)

	// If access Token is invalid we clear the cookie and throw error
	if err != nil {
		c.ClearCookie()
		config.Logger.Error("Inavlid " + ACCESS_TOKEN + " Token: " + err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	config.Logger.Info(ACCESS_TOKEN + " Token is valid")

	// Get Token type from Token
	tokenType, err := GetTypeFromValidatedToken(validatedToken)
	if err != nil {
		config.Logger.Error(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Verifying Token type (it must be access)
	if err := VerifyTokenType(ACCESS_TOKEN, tokenType); err != nil {
		config.Logger.Error(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	config.Logger.Info("Token type is correct")

	// Extracting User ID from valid access Token
	userId, err := GetIdFromValidatedToken(validatedToken)
	if err != nil {
		config.Logger.Error(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	config.Logger.Info("User ID from " + ACCESS_TOKEN + " Token is: " + fmt.Sprint(userId))

	// Search User by Id in DB
	DBUser, err := FindUserById(int(userId))
	if err != nil {
		c.ClearCookie()
		config.Logger.Error("No User with such id: " + fmt.Sprint(userId))
		return fiber.NewError(fiber.StatusBadRequest, "No such User")
	}
	config.Logger.Info("User found successfully")

	return c.JSON(ConvertUserToResponse(accessToken, "", &DBUser))
}

// Refresh Access Token
// @Summary Verifying refresh Token and returning Access
// @Description Verifying access, extract sub and returns Token status. Clears the Cookies, if there any error
// @Tags authentication
// @Produce json
// @Success 200 {object} models.UserResponse "Refresh Token is Valid"
// @Failure 400 {object} models.ErrorResponse "Invalid access Token or No such User"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /auth/refresh [get]
func Refresh(c *fiber.Ctx) error {
	var refreshToken string = c.Cookies(REFRESH_TOKEN)

	// Refresh Token Verifying
	validatedToken, err := TokenVerify(refreshToken)

	// If refresh Token is invalid we clear the cookie and throw error
	if err != nil {
		c.ClearCookie()
		config.Logger.Error("Inavlid " + REFRESH_TOKEN + " Token: " + err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	config.Logger.Info(REFRESH_TOKEN + " Token is valid")

	// Get Token type from Token
	tokenType, err := GetTypeFromValidatedToken(validatedToken)
	if err != nil {
		config.Logger.Error(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Verifying Token type (it must be refresh)
	if err := VerifyTokenType(REFRESH_TOKEN, tokenType); err != nil {
		config.Logger.Error(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	config.Logger.Info("Token type is correct")

	// Extracting User ID from valid refresh Token
	userId, err := GetIdFromValidatedToken(validatedToken)
	if err != nil {
		config.Logger.Error(err.Error())
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	config.Logger.Info("User ID from " + REFRESH_TOKEN + " Token is: " + fmt.Sprint(userId))

	// Search User by Id in DB
	DBUser, err := FindUserById(int(userId))
	if err != nil {
		c.ClearCookie()
		config.Logger.Error("No User with such id: " + fmt.Sprint(userId))
		return fiber.NewError(fiber.StatusBadRequest, "No such User")
	}
	config.Logger.Info("User found successfully")

	// Creating Access Token
	accessToken, err := CreateToken(ACCESS_TOKEN, &DBUser)
	if err != nil {
		config.Logger.Error("Problem with creating Access JWT Token" + err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Problem with creating Access JWT Token")
	}
	c.Cookie(&fiber.Cookie{
		Name:  ACCESS_TOKEN,
		Value: accessToken,
	})
	config.Logger.Info("Access JWT created successfully")

	return c.JSON(ConvertUserToResponse(accessToken, refreshToken, &DBUser))
}
