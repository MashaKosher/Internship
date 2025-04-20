package handler

import (
	routes "adminservice/internal/handler/internal"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func Handlers(r *chi.Mux) {
	// Middlewares
	middleWares(r)

	// Swagger route initialize
	routes.SwaggerRoutes(r)

	// Auth routes initialize
	routes.AdminRoutes(r)
}

func middleWares(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
}
