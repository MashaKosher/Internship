package http

import (
	customMiddleware "coreservice/internal/controller/http/middleware"
	routes "coreservice/internal/controller/http/v1"
	"coreservice/internal/di"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine, deps di.Container) {

	router.Use(customMiddleware.PrometheusMiddleware())
	router.Use(customMiddleware.TracingMiddleware())

	go func() {
		for {
			customMiddleware.UpdateMetrics()
			time.Sleep(5 * time.Second) // Обновление каждые 5 секунд
		}
	}()

	routes.IntiMetricsRoutes(router)

	// Swagger Routes
	routes.SwaggerRoutes(router)

	routes.SearchRoutes(router, deps)

	// Auth middleware for all routes
	router.Use(customMiddleware.AuthMiddleWare(deps.Config, deps.Logger, deps.DB, deps.Elastic, deps.Bus))

	routes.UserRoutes(router, deps)

	routes.SeasonRoutes(router, deps)

	routes.TokenRoutes(router, deps)

	routes.DailyTasksRoutes(router, deps)

}
