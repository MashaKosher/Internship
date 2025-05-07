package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func IntiMetricsRoutes(router *gin.Engine) {
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
