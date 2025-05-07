package controller

import (
	customMiddleware "gameservice/internal/controller/middlewares"
	routes "gameservice/internal/controller/v1"
	"gameservice/internal/di"
	"time"

	"github.com/labstack/echo/v4"
)

func NewRouter(e *echo.Echo, deps di.Container) {

	e.Use(customMiddleware.MetricsMiddleware)
	e.Use(customMiddleware.TracingMiddleware)
	go func() {
		for {
			customMiddleware.UpdateMetrics()
			time.Sleep(5 * time.Second)
		}
	}()

	// Добавляем хэндлер для метрик Prometheus

	routes.IntiMetricsRoutes(e)

	// Swagger route initialize
	routes.InitSwaggerRoutes(e)
	routes.InitWSRoutes(e, deps)
	routes.InitGameRoutes(e, deps)

}
