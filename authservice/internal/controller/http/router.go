package http

import (
	routes "authservice/internal/controller/http/v1"
	"authservice/internal/di"

	// "net/http"

	"github.com/gofiber/fiber/v2"
	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Создаем метрики
// var (
// 	httpRequestsTotal = prometheus.NewCounterVec(
// 		prometheus.CounterOpts{
// 			Name: "http_requests_total",
// 			Help: "Total HTTP requests",
// 		},
// 		[]string{"method", "endpoint"},
// 	)
// 	httpRequestDuration = prometheus.NewHistogramVec(
// 		prometheus.HistogramOpts{
// 			Name:    "http_request_duration_seconds",
// 			Help:    "HTTP request duration",
// 			Buckets: prometheus.DefBuckets,
// 		},
// 		[]string{"method", "endpoint"},
// 	)
// )

func NewRouter(app *fiber.App, deps di.Container) {
	// app.Use(custoMiddleware.MetricsMiddleware)

	// prometheus.MustRegister(httpRequestsTotal)
	// prometheus.MustRegister(httpRequestDuration)

	// app.Get("/metrics", func(c *fiber.Ctx) error {
	// 	promhttp.Handler().ServeHTTP(c.Context().Response().BodyWriter(), &http.Request{})
	// 	return nil
	// })

	// Swagger route initialize
	routes.InitSwaggerRoutes(app)

	// Auth routes initialize
	routes.InitAuthRoutes(app, deps)
}
