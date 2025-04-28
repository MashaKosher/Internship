package main

import (
	"gameservice/internal/app"
	"gameservice/internal/config"

	_ "gameservice/docs" // подключение сгенерированной документации Swagger
)

// @title						Game service
// @version						1.0
// @description					Auth server API
// @host						localhost:8005
// @BasePath					/
// @securityDefinitions.basic	BasicAuth
// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	config.Load()

	app.Run()
}
