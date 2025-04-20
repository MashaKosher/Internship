package internal

import (
	"coreservice/internal/service"

	"github.com/gin-gonic/gin"
)

func TokenRoutes(app *gin.Engine) {
	coreRoutes := app.Group("/check-token")
	coreRoutes.GET("/", service.CheckToken)
}
