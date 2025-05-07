package http

import (
	customMiddleware "adminservice/internal/controller/http/middlewares"
	routes "adminservice/internal/controller/http/v1"
	"adminservice/internal/di"
	"time"

	"github.com/go-chi/chi/v5"
)

func NewRouter(r *chi.Mux, deps di.Container) {

	// Middlewares
	middleWares(r)

	go func() {
		for {
			customMiddleware.UpdateMetrics()
			time.Sleep(5 * time.Second) // Обновление каждые 5 секунд
		}
	}()

	routes.IntiMetricsRoutes(r)

	// Swagger route initialize
	routes.InitSwaggerRoutes(r)

	// Auth routes initialize
	routes.InitAdminRoutes(r, deps)
}

func middleWares(r *chi.Mux) {
	r.Use(customMiddleware.MetricsMiddleware)
	r.Use(customMiddleware.TracingMiddleware)
}
