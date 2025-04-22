package internal

import (
	"github.com/go-chi/chi/v5"

	controllers "adminservice/internal/service"
)

func SwaggerRoutes(r *chi.Mux) {
	r.Group(func(r chi.Router) {
		r.Get("/swagger/*", controllers.Swagger())
	})
}
