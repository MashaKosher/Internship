package http

import (
	customMiddleware "authservice/internal/controller/http/middlewares"
	routes "authservice/internal/controller/http/v1"
	"authservice/internal/di"

	"time"

	"github.com/gofiber/fiber/v2"
)

func NewRouter(app *fiber.App, deps di.Container) {

	app.Use(customMiddleware.MetricsMiddleware)
	app.Use(customMiddleware.TracingMiddleware)

	go func() {
		for {
			customMiddleware.UpdateMetrics()
			time.Sleep(5 * time.Second) // Обновление каждые 5 секунд
		}
	}()

	// Эндпоинт для Prometheus (исправленный)
	routes.InitMetricsRoutes(app)

	// Swagger route initialize
	routes.InitSwaggerRoutes(app)

	// Auth routes initialize
	routes.InitAuthRoutes(app, deps)
}
