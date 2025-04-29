package app

import (
	v1 "adminservice/internal/controller/http"
	"adminservice/internal/di"

	"net/http"

	"github.com/go-chi/chi/v5"
)

func Run(deps di.Container) http.Handler {

	r := chi.NewRouter()

	v1.NewRouter(r, deps)

	return r
}
