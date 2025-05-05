package http

import (
	"coreservice/internal/controller/http/middleware"
	routes "coreservice/internal/controller/http/v1"
	"coreservice/internal/di"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total_core",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path"},
	)
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds_core",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{0.1, 0.3, 0.5, 1, 2, 5},
		},
		[]string{"method", "path"},
	)
)

func init() {
	// Регистрируем метрики
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

func NewRouter(router *gin.Engine, deps di.Container) {

	router.Use(func(c *gin.Context) {
		// Пропускаем метрики, чтобы не учитывать их в статистике
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		// Замер времени выполнения
		timer := prometheus.NewTimer(httpRequestDuration.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		))
		defer timer.ObserveDuration()

		// Увеличиваем счетчик запросов
		httpRequestsTotal.WithLabelValues(
			c.Request.Method,
			c.FullPath(),
		).Inc()

		c.Next()
	})

	// Добавляем endpoint для Prometheus
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger Routes
	routes.SwaggerRoutes(router)

	routes.SearchRoutes(router, deps)

	// Auth middleware for all routes
	router.Use(middleware.AuthMiddleWare(deps.Config, deps.Logger, deps.DB, deps.Elastic, deps.Bus))

	routes.UserRoutes(router, deps)

	routes.SeasonRoutes(router, deps)

	routes.TokenRoutes(router, deps)

	routes.DailyTasksRoutes(router, deps)

}
