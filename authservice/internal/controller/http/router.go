package http

import (
	routes "authservice/internal/controller/http/v1"

	"github.com/gofiber/fiber/v2"

	"authservice/internal/usecase"
)

func NewRouter(app *fiber.App, authUsecase usecase.Auth) {
	// Swagger route initialize
	routes.InitSwaggerRoutes(app)

	// Auth routes initialize
	routes.InitAuthRoutes(app, authUsecase)
}
