package v1

import (
	"authservice/internal/di"
	"authservice/internal/entity"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type authRoutes struct {
	u di.AuthService
	v di.ValidatorType
	l di.LoggerType
}

func InitAuthRoutes(app *fiber.App, deps di.Container) {

	r := &authRoutes{
		u: deps.Services.Auth,
		v: deps.Validator,
		l: deps.Logger,
	}

	api := app.Group("/auth")
	{
		api.Post("/login", r.login)

		check := api.Group("/check")
		{
			check.Get("/access", r.checkAccessToken)
			check.Get("/refresh", r.checkRefreshToken)
			check.Get("/", r.checkTokens)

		}

		signUp := api.Group("/sign-up")
		{
			signUp.Post("/user", r.userSignUp)
			signUp.Post("/admin", r.adminSignUp)

		}
		api.Post("/change-password", r.changePassword)

		api.Delete("/delete", r.deleteUser)
	}
}

// @Summary      User login
// @Description  Authenticates user and returns JWT tokens in cookies and response body
// @Tags         Login
// @Accept       json
// @Produce      json
// @Param		 entity.UserInDTO	body  entity.UserInDTO	true	"Login request body"
// @Success      200 {object} entity.UserOutDTO "Successfully logged in"
// @Header       200 {string} Set-Cookie "access_token=JWT_TOKEN; Path=/; HttpOnly"
// @Header       200 {string} Set-Cookie "refresh_token=JWT_TOKEN; Path=/; HttpOnly"
// @Failure      400 {object} entity.Error "Invalid input data"
// @Failure      404  {object}  entity.Error     "Not Found - User not found"
// @Failure      500 {object} entity.Error "Internal server error"
// @Router       /auth/login [post]
func (r *authRoutes) login(c *fiber.Ctx) error {
	var user entity.UserInDTO
	c.BodyParser(&user)
	r.l.Info(fmt.Sprintf("User with UserName %s is trying to log in", user.Username))

	// Veryfing income credentials
	if err := r.v.Struct(user); err != nil {
		r.l.Error("Body validation error: " + err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid data")
	}
	r.l.Info("User credentials verified successfully")

	outUser, err := r.u.Login(user)
	if err != nil {
		c.ClearCookie()
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:  string(di.ACCESS_TOKEN),
		Value: outUser.AccessToken,
	})
	r.l.Info("Access JWT set successfully")

	c.Cookie(&fiber.Cookie{
		Name:  string(di.REFRESH_TOKEN),
		Value: outUser.RefreshToken,
	})
	r.l.Info("Refresh JWT set successfully")

	return c.JSON(outUser)
}

// @Summary      Verify refresh token
// @Description  Verifies JWT refresh token from cookies and returns user data if valid. Clears cookies on any error.
// @Tags         Check Token
// @Produce      json
// @Success      200 {object} entity.UserInDTO "Refresh token is valid"
// @Failure      400 {object} entity.Error "Bad Request - Missing or empty tokens"
// @Failure      401 {object} entity.Error "Unauthorized - Invalid or expired token"
// @Failure      403 {object} entity.Error "Forbidden - Token validation failed"
// @Failure      404 {object} entity.Error "Not Found - User not found"
// @Failure      500 {object} entity.Error "Internal server error"
// @Router       /auth/check/access [get]
func (r *authRoutes) checkAccessToken(c *fiber.Ctx) error {
	var accessToken string = c.Cookies(string(di.ACCESS_TOKEN))

	if len(accessToken) == 0 {
		r.l.Error("Access token field is empty")
		return fiber.NewError(fiber.StatusBadRequest, "Access token field is empty")
	}

	outUser, err := r.u.CheckAccessToken(accessToken)
	if err != nil {
		c.ClearCookie()
		return err
	}
	return c.JSON(outUser)
}

// @Summary      Verify refresh token
// @Description  Verifies JWT refresh token from cookies and returns user data if valid. Clears cookies on any error.
// @Tags         Check Token
// @Produce      json
// @Success      200 {object} entity.UserInDTO "Refresh token is valid"
// @Failure      400  {object}  entity.Error      "Bad Request - Missing or empty tokens"
// @Failure      401 {object} entity.Error "Unauthorized - Invalid or expired token"
// @Failure      403 {object} entity.Error "Forbidden - Token validation failed"
// @Failure      404  {object}  entity.Error      "Not Found - User not found"
// @Failure      500 {object} entity.Error "Internal server error"
// @Router       /auth/check/refresh [get]
func (r *authRoutes) checkRefreshToken(c *fiber.Ctx) error {
	var refreshToken string = c.Cookies(string(di.REFRESH_TOKEN))

	if len(refreshToken) == 0 {
		r.l.Error("Refresh token field is empty")
		return fiber.NewError(fiber.StatusBadRequest, "Refresh token field is empty")
	}

	outUser, err := r.u.CheckRefreshToken(refreshToken)
	if err != nil {
		c.ClearCookie()
		return err
	}
	return c.JSON(outUser)
}

// @Summary      Verify both tokens
// @Description  Verifies both access and refresh JWT tokens from cookies. Returns user data if refresh token is valid. Clears cookies on any error.
// @Tags         Check Token
// @Produce      json
// @Success      200  {object}  entity.UserInDTO  "Tokens are valid, returns user data"
// @Failure      400  {object}  entity.Error      "Bad Request - Missing or empty tokens"
// @Failure      401  {object}  entity.Error      "Unauthorized - Invalid or expired tokens"
// @Failure      403  {object}  entity.Error      "Forbidden - Token validation failed"
// @Failure      404  {object}  entity.Error      "Not Found - User not found"
// @Failure      500  {object}  entity.Error      "Internal server error"
// @Router       /auth/check [get]
func (r *authRoutes) checkTokens(c *fiber.Ctx) error {
	var accessToken string = c.Cookies(string(di.ACCESS_TOKEN))
	if len(accessToken) == 0 {
		r.l.Error("Access token field is empty")
		return fiber.NewError(fiber.StatusBadRequest, "Access token field is empty")
	}

	var refreshToken string = c.Cookies(string(di.REFRESH_TOKEN))
	if len(refreshToken) == 0 {
		r.l.Error("Refresh token field is empty")
		return fiber.NewError(fiber.StatusBadRequest, "Refresh token field is empty")
	}

	outUser, err := r.u.CheckTokens(accessToken, refreshToken)
	if err != nil {
		c.ClearCookie()
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:  string(di.ACCESS_TOKEN),
		Value: outUser.AccessToken,
	})

	return c.JSON(outUser)
}

// @Summary      Register new user
// @Description  Creates a new user account with default User role
// @Tags         Sign Up
// @Accept       json
// @Produce      json
// @Param        request body entity.UserSignUpInDTO true "User registration data"
// @Success      201 {object} entity.UserInDTO "Successfully registered"
// @Failure      400 {object} entity.Error "Invalid input data"
// @Failure      409 {object} entity.Error "Conflict - Username already exists"
// @Failure      500 {object} entity.Error "Internal server error"
// @Router       /auth/sign-up/user [post]
func (r *authRoutes) userSignUp(c *fiber.Ctx) error {
	var user entity.UserSignUpInDTO
	c.BodyParser(&user)
	r.l.Info(fmt.Sprintf("User with UserName %s is trying to log in", user.Username))

	// Veryfing income credentials
	if err := r.v.Struct(user); err != nil {
		r.l.Error("Body validation error: " + err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid data")
	}
	r.l.Info("User credentials verified successfully")

	outUser, err := r.u.UserSignUp(user.User, user.ReferalID)
	if err != nil {
		c.ClearCookie()
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:  string(di.ACCESS_TOKEN),
		Value: outUser.AccessToken,
	})
	r.l.Info("Access JWT set successfully")

	c.Cookie(&fiber.Cookie{
		Name:  string(di.REFRESH_TOKEN),
		Value: outUser.RefreshToken,
	})
	r.l.Info("Refresh JWT set successfully")

	return c.JSON(outUser)
}

// @Summary      Register new admin
// @Description  Creates a new user account with Admin privileges (requires special permissions)
// @Tags         Sign Up
// @Accept       json
// @Produce      json
// @Param        request body entity.UserSignUpInDTO true "Admin registration data"
// @Success      201 {object} entity.UserInDTO  "Admin successfully registered"
// @Failure      400 {object} entity.Error "Invalid input data"
// @Failure      401 {object} entity.Error "Unauthorized - Only existing admins can create new admins"
// @Failure      409 {object} entity.Error "Conflict - Username already exists"
// @Failure      500 {object} entity.Error "Internal server error"
// @Router       /auth/sign-up/admin [post]}
func (r *authRoutes) adminSignUp(c *fiber.Ctx) error {
	var user entity.UserSignUpInDTO
	c.BodyParser(&user)
	r.l.Info(fmt.Sprintf("User with UserName %s is trying to log in", user.Username))

	// Veryfing income credentials
	if err := r.v.Struct(user); err != nil {
		r.l.Error("Body validation error: " + err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid data")
	}
	r.l.Info("User credentials verified successfully")

	outUser, err := r.u.AdminSignUp(user.User, user.ReferalID)
	if err != nil {
		c.ClearCookie()
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:  string(di.ACCESS_TOKEN),
		Value: outUser.AccessToken,
	})
	r.l.Info("Access JWT set successfully")

	c.Cookie(&fiber.Cookie{
		Name:  string(di.REFRESH_TOKEN),
		Value: outUser.RefreshToken,
	})
	r.l.Info("Refresh JWT set successfully")

	return c.JSON(outUser)
}

// @Summary      Change user password
// @Description  Changes password for authenticated user. Requires valid access token in cookies.
// @Tags         Change Password
// @Accept       json
// @Produce      json
// @Param        request body entity.Password true "New password data"
// @Success      200  {object}  entity.UserInDTO  "Tokens are valid, returns user data"
// @Failure      400  {object}  entity.Error      "Bad Request - Missing or empty tokens"
// @Failure      401  {object}  entity.Error      "Unauthorized - Invalid or expired tokens"
// @Failure      403  {object}  entity.Error      "Forbidden - Token validation failed"
// @Failure      404  {object}  entity.Error      "Not Found - User not found"
// @Failure      500  {object}  entity.Error      "Internal server error"
// @Router       /auth/change-password [post]
func (r *authRoutes) changePassword(c *fiber.Ctx) error {
	var newPassword entity.Password
	c.BodyParser(&newPassword)

	// Veryfing income credentials
	if err := r.v.Struct(newPassword); err != nil {
		r.l.Error("Body validation error: " + err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid data")
	}

	var accessToken string = c.Cookies(string(di.ACCESS_TOKEN))

	outUser, err := r.u.ChangePassword(newPassword, accessToken)
	if err != nil {
		c.ClearCookie()
		return err
	}
	return c.JSON(outUser)

}

// @Summary      Delete user account
// @Description  Permanently deletes user account after validating both access and refresh tokens. Clears all auth cookies on any error.
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200  {object}  entity.UserInDTO  "Tokens are valid, returns user data"
// @Failure      400  {object}  entity.Error      "Bad Request - Missing or empty tokens"
// @Failure      401  {object}  entity.Error      "Unauthorized - Invalid or expired tokens"
// @Failure      403  {object}  entity.Error      "Forbidden - Token validation failed"
// @Failure      404  {object}  entity.Error      "Not Found - User not found"
// @Failure      500  {object}  entity.Error      "Internal server error"
// @Router       /auth/delete [delete]
func (r *authRoutes) deleteUser(c *fiber.Ctx) error {
	var accessToken string = c.Cookies(string(di.ACCESS_TOKEN))
	if len(accessToken) == 0 {
		r.l.Error("Access token field is empty")
		return fiber.NewError(fiber.StatusBadRequest, "Access token field is empty")
	}

	var refreshToken string = c.Cookies(string(di.REFRESH_TOKEN))
	if len(refreshToken) == 0 {
		r.l.Error("Refresh token field is empty")
		return fiber.NewError(fiber.StatusBadRequest, "Refresh token field is empty")
	}

	outUser, err := r.u.CheckTokens(accessToken, refreshToken)
	if err != nil {
		c.ClearCookie()
		return err
	}

	if err = r.u.DeleteUser(outUser.UserID); err != nil {
		r.l.Error("Some problems with user delete")
		return err
	}

	return c.JSON(outUser)
}
