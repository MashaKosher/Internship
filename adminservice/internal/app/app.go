package app

import (
	v1 "adminservice/internal/controller/http"
	db "adminservice/pkg/client/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Run() http.Handler {
	r := chi.NewRouter()

	_ = db.ConncetDB()
	v1.NewRouter(r)

	return r
}
