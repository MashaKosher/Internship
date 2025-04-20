package internal

import (
	"coreservice/internal/service"

	"github.com/gin-gonic/gin"
)

func ElasticRoutes(router *gin.Engine) {
	coreRoutes := router.Group("/elastic")
	coreRoutes.GET("/create-index", service.UserBuildSearchIndex)
	coreRoutes.POST("/strict", service.SearchElasticByNameStrict)
	coreRoutes.POST("/wildcard", service.SearchElasticByNameWildcard)
	coreRoutes.POST("/fuzzy", service.SearchElasticByNameFuzzy)
}
