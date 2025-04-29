package main

import (
	_ "authservice/docs"
	"authservice/internal/app"

	"authservice/internal/config"
)

// @title						Auth service
// @version						1.0
// @description					Auth server API
// @host						localhost:8080
// @BasePath					/
// @securityDefinitions.basic	BasicAuth
// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	cfg := config.MustParseConfig()
	app.Run(cfg)
}
