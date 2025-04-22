package http

import (
	routes "authservice/internal/controller/http/v1"

	"github.com/gofiber/fiber/v2"
)

func NewRouter(app *fiber.App) {
	// Swagger route initialize
	routes.SwaggerRoutes(app)

	// Auth routes initialize
	routes.AuthRoutes(app)
}
