package http

import (
	routes "adminservice/internal/controller/http/v1"
	"adminservice/internal/di"

	custoMiddleware "adminservice/internal/controller/http/middlewares"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(r *chi.Mux, deps di.Container) {
	// Middlewares
	middleWares(r)

	routes.IntiMetricsRoutes(r)

	// Swagger route initialize
	routes.InitSwaggerRoutes(r)

	// Auth routes initialize
	routes.InitAdminRoutes(r, deps)
}

func middleWares(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(custoMiddleware.MetricsMiddleware)
}
