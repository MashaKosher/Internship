package v1

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger" // для подключения Swagger UI
)

func InitSwaggerRoutes(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
