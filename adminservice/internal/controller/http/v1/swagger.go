package v1

import (
	"github.com/go-chi/chi/v5"

	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func InitSwaggerRoutes(r *chi.Mux) {
	r.Group(func(r chi.Router) {
		r.Get("/swagger/*", swagger())
	})
}

func swagger() http.HandlerFunc {

	url := httpSwagger.URL("http://localhost:8004/swagger/doc.json")

	return httpSwagger.Handler(url)

}
