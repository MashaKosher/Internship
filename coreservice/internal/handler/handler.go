package handler

import (
	routes "coreservice/internal/handler/internal"
	"coreservice/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Handlers(router *gin.Engine) {
	// Swagger Routes
	routes.SwaggerRoutes(router)

	routes.ElasticRoutes(router)

	// Auth middleware for all routes
	router.Use(middleware.AuthMiddleWare())

	routes.UserRoutes(router)

	routes.TokenRoutes(router)

	routes.SeasonRoutes(router)

	routes.DailyTasksRoutes(router)

}
