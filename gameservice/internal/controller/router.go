package controller

import (
	routes "gameservice/internal/controller/v1"
	"gameservice/internal/di"

	"github.com/labstack/echo/v4"
)

func NewRouter(e *echo.Echo, deps di.Container) {

	// Swagger route initialize
	routes.InitSwaggerRoutes(e)
	routes.InitWSRoutes(e, deps)
	routes.InitGameRoutes(e, deps)

}
