package v1

import (
	"github.com/go-chi/chi/v5"

	controllers "adminservice/internal/service"
)

func InitSwaggerRoutes(r *chi.Mux) {
	r.Group(func(r chi.Router) {
		r.Get("/swagger/*", controllers.Swagger())
	})
}
