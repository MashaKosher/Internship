package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func IntiMetricsRoutes(e *echo.Echo) {
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}
