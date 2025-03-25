package internal

import (
	controllers "authservice/controllers/auth"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	api := app.Group("/auth")
	api.Post("/signup", controllers.SignUp)
	api.Post("/login", controllers.Login)
	api.Get("/check-token", controllers.CheckToken)
	api.Get("/refresh", controllers.Refresh)
}
