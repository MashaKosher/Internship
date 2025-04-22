package internal

import (
	"authservice/internal/service"
	"authservice/pkg/logger"
	"authservice/pkg/tokens"
	"fmt"

	"authservice/internal/usecase"

	dto "authservice/internal/usecase/auth"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type authRoutes struct {
	u usecase.Auth
	v *validator.Validate
}

func InitAuthRoutes(app *fiber.App, authUsecase usecase.Auth) {

	r := &authRoutes{u: authUsecase, v: validator.New(validator.WithRequiredStructEnabled())}

	api := app.Group("/auth")
	{
		api.Post("/login", r.login)
		// ////////////////////////////////
		api.Get("/check-token", service.CheckToken)
		api.Get("/refresh", service.Refresh)
		signUp := api.Group("/sign-up")
		{
			signUp.Post("/user", service.UserSignUp)
			signUp.Post("/admin", service.AdminSignUp)
		}
		api.Post("/change-password", service.ChangePassword)
	}
}

// @Summary      User login
// @Description  Authenticates user and returns JWT tokens in cookies and response body
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param		 auth.UserInDTO	body  auth.UserInDTO	true	"Login request body"
// @Success      200 {object} auth.LoginOutDTO "Successfully logged in"
// @Header       200 {string} Set-Cookie "access_token=JWT_TOKEN; Path=/; HttpOnly"
// @Header       200 {string} Set-Cookie "refresh_token=JWT_TOKEN; Path=/; HttpOnly"
// @Failure      400 {object} entity.Error "Invalid input data"
// @Failure      401 {object} entity.Error "Unauthorized"
// @Failure      500 {object} entity.Error "Internal server error"
// @Router       /auth/login [post]
func (r *authRoutes) login(c *fiber.Ctx) error {
	var user dto.UserInDTO
	c.BodyParser(&user)
	logger.Logger.Info(fmt.Sprintf("User with UserName %s is trying to log in", user.Username))

	// Veryfing income credentials
	if err := r.v.Struct(user); err != nil {
		logger.Logger.Error("Body validation error: " + err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid data")
	}
	logger.Logger.Info("User credentials verified successfully")

	outUser, err := r.u.Login(user)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:  tokens.ACCESS_TOKEN,
		Value: outUser.AccessToken,
	})
	logger.Logger.Info("Access JWT created successfully")

	c.Cookie(&fiber.Cookie{
		Name:  tokens.REFRESH_TOKEN,
		Value: outUser.RefreshToken,
	})
	logger.Logger.Info("Refresh JWT created successfully")

	return c.JSON(outUser)
}
