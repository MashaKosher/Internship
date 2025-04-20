package internal

import (
	"coreservice/internal/service"

	"github.com/gin-gonic/gin"
)

func DailyTasksRoutes(router *gin.Engine) {

	router.GET("/daily-task", service.DailyTask)
	//		coreRoutes := router.Group("/elastic")
	//		coreRoutes.GET("/create-index", service.UserBuildSearchIndex)
	//		coreRoutes.POST("/strict", service.SearchElasticByNameStrict)
	//		coreRoutes.POST("/wildcard", service.SearchElasticByNameWildcard)
	//		coreRoutes.POST("/fuzzy", service.SearchElasticByNameFuzzy)
	//	}
}
