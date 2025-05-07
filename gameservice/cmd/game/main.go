package main

import (
	"gameservice/internal/app"
	"gameservice/internal/config"

	_ "gameservice/docs"
)

// @title						Game service
// @version						1.0
// @description					Game server API
// @host						localhost:8005
// @BasePath					/
// @securityDefinitions.basic	BasicAuth
// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	cfg := config.MustParseConfig()
	app.Run(cfg)
}
