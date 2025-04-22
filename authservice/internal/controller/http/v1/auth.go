package internal

import (
	"authservice/internal/service"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {

	api := app.Group("/auth")
	{
		api.Post("/login", service.Login)
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
