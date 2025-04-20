package internal

import (
	"coreservice/internal/service"

	"github.com/gin-gonic/gin"
)

func SeasonRoutes(router *gin.Engine) {
	seasonRoutes := router.Group("/seasons")
	seasonRoutes.GET("/", service.Seasons)
	seasonRoutes.GET("/current", service.CurrentSeason)
	seasonRoutes.GET("/planned", service.PlannedSeason)
	seasonRoutes.GET("/:id/leader-board", service.SeasonLeaderBoard)
	seasonRoutes.GET("/:id", service.SeasonInfo)
}
