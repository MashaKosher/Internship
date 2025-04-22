package usecase

import "authservice/internal/usecase/auth"

type (
	// Translation -.
	Auth interface {
		Login(user auth.UserInDTO) (auth.LoginOutDTO, error)
		// CheckToken(c *fiber.Ctx) error
		// Refresh(c *fiber.Ctx) error
		// UserSignUp(c *fiber.Ctx) error
		// AdminSignUp(c *fiber.Ctx) error
		// ChangePassword(c *fiber.Ctx) error

	}
)
