package service

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func Swagger() http.HandlerFunc {

	url := httpSwagger.URL("http://localhost:8004/swagger/doc.json")

	return httpSwagger.Handler(url)

}
