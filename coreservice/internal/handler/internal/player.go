package internal

import (
	"coreservice/internal/service"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	coreRoutes := router.Group("/user")
	coreRoutes.POST("/deposit", service.MakeDeposit)
	coreRoutes.GET("/info", service.UserInfo)
}
