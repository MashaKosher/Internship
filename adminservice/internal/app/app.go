package app

import (
	"adminservice/internal/handler"
	db "adminservice/pkg/client/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Run() http.Handler {
	r := chi.NewRouter()

	_ = db.ConncetDB()
	handler.Handlers(r)

	return r
}
