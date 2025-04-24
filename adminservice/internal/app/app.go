package app

import (
	"adminservice/internal/db"
	"adminservice/internal/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Run() http.Handler {
	r := chi.NewRouter()

	db.ConncetDB()
	handler.Handlers(r)

	return r
}
