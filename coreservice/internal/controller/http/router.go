package http

import (
	"coreservice/internal/controller/http/middleware"
	routes "coreservice/internal/controller/http/v1"
	"coreservice/internal/di"

	"github.com/gin-gonic/gin"
)

func NewRouter(router *gin.Engine, deps di.Container) {
	// Swagger Routes
	routes.SwaggerRoutes(router)

	routes.SearchRoutes(router, deps)

	// Auth middleware for all routes
	router.Use(middleware.AuthMiddleWare(deps.Logger))

	routes.UserRoutes(router, deps)

	routes.SeasonRoutes(router, deps)

	routes.TokenRoutes(router, deps)

	routes.DailyTasksRoutes(router, deps)

}
