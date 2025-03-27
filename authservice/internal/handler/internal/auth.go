package internal

import (
	"authservice/internal/service"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	api := app.Group("/auth")
	api.Post("/signup", service.SignUp)
	api.Post("/login", service.Login)
	api.Get("/check-token", service.CheckToken)
	api.Get("/refresh", service.Refresh)
}
