package auth

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"

	models "authservice/models"

	config "authservice/config"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// SignUp обрабатывает запросы на регистрацию пользователя.
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

	if err := validate.Struct(user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid data")
	}

	// HASHING PASSWORD
	hashed, err := HashPassword(user.Password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Problems with hashing password")
	}
	user.Password = hashed

	// Adding User to DB
	if err := CreateUser(&user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Problem with adding user to DB")
	}

	// GENERATING TOKEN
	config.Logger.Info("User created: " + fmt.Sprint(user))
	token, err := GenerateToken(&user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Problem with generating JWT")
	}

	c.Cookie(&fiber.Cookie{
		Name:  "jwt",
		Value: token,
	})

	return c.JSON(ConvertUserToResponse(token, &user))

}

// Login обрабатывает запросы на вход пользователя.
// @Summary User login
// @Description Returns a message indicating the login endpoint
// @Tags authentication
// @Accept json
// @Produce json
// @Param models.User body models.User true "Login request body"
// @Success 200 {object} models.UserResponse "User successfully logged"
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
	var user models.User
	c.BodyParser(&user)

	if err := validate.Struct(user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid data")
	}

	// IF USER WITH SUCH NAME EXISTS
	DBUser, err := FindUser(&user)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Username")
	}

	// COMPARING HASHED PASSWORD FROM DB WITH RAW PASSWORD FROM JSON
	if err := ValidatePassword(DBUser.Password, user.Password); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Password")
	}

	// GENERATING TOKEN
	token, err := GenerateToken(&DBUser)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Problem with generating JWT")
	}

	c.Cookie(&fiber.Cookie{
		Name:  "jwt",
		Value: token,
	})

	return c.JSON(ConvertUserToResponse(token, &DBUser))
}

// CheckToken проверяет JWT токен, переданный в куках.
// @Summary Проверка токена JWT
// @Description Проверяет токен JWT, извлекает userId и возвращает статус токена. Если токен недействителен, очищает куку.
// @Tags authentication
// @Produce json
// @Success 200
// @Router /auth/check-token [get]
func CheckToken(c *fiber.Ctx) error {
	var token string = c.Cookies("jwt")

	validatedToken, err := TokenValidation(token)

	// if token is invalid we clear the cookie
	if err != nil {
		config.Logger.Error("Inavlid Token")
		c.ClearCookie()
		return fiber.NewError(fiber.StatusBadRequest, "Unexpceted signing method")
	}

	unvalidatedUserId := validatedToken.Claims.(jwt.MapClaims)["userId"]

	userId, ok := unvalidatedUserId.(float64)
	if !ok {
		config.Logger.Error("Invalid ID value:" + fmt.Sprint(unvalidatedUserId))
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Id")
	}

	if err := FindUserById(int(userId)); err != nil {
		config.Logger.Error("No User with such id: " + fmt.Sprint(userId))
		c.ClearCookie()
		return fiber.NewError(fiber.StatusBadRequest, "No such User")
	}

	return c.JSON(fiber.Map{
		"token status": userId,
	})
}
