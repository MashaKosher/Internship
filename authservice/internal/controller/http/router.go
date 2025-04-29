package http

import (
	routes "authservice/internal/controller/http/v1"
	"authservice/internal/di"

	"github.com/gofiber/fiber/v2"
)

func NewRouter(app *fiber.App, deps di.Container) {
	// Swagger route initialize
	routes.InitSwaggerRoutes(app)

	// Auth routes initialize
	routes.InitAuthRoutes(app, deps)
}
