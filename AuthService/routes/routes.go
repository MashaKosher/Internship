package routes

import (
	routes "authservice/routes/internal"

	"github.com/gofiber/fiber/v2"
)

func Handlers(app *fiber.App) {
	// Swagger route initialize
	routes.SwaggerRoutes(app)

	// Auth routes initialize
	routes.AuthRoutes(app)
}
