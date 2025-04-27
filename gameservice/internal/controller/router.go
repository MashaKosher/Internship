package controller

import (
	routes "gameservice/internal/controller/v1"
	"gameservice/internal/usecase"

	"github.com/labstack/echo/v4"
)

func NewRouter(e *echo.Echo, gameUseCase usecase.Game) {

	// Swagger route initialize
	routes.InitSwaggerRoutes(e)
	routes.InitGameRoutes(e, gameUseCase)

}
