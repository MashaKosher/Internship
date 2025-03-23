package internal

import (
	controllers "authservice/controllers/auth"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {

	app.Get("/", hello)

	api := app.Group("/auth")
	api.Post("/signup", controllers.SignUp)
	api.Post("/login", controllers.Login)
	api.Get("/check-token", controllers.CheckToken)
}

// hello возвращает приветственное сообщение.
// @Summary Get greeting message
// @Description Returns a simple greeting message
// @Tags greeting
// @Produce json
// @Success 200
// @Router / [get]
func hello(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"Message": "Hello",
	})
}
