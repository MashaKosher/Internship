package main

import (
	_ "authservice/docs"
	"authservice/internal/adapter/kafka/consumers"
	"authservice/internal/handler"
	"fmt"

	"authservice/internal/config"
	"authservice/internal/db"
	"authservice/internal/keys"

	"authservice/internal/logger"

	"github.com/gofiber/fiber/v2"
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

	config.Load()

	fmt.Println("File Log Name: " + config.AppConfig.Logger.FileName)

	// Creating Log File
	logFile := logger.CreateLogger()
	defer logFile.Close()
	defer logger.Logger.Sync()
	defer logger.Logger.Info("Program end")

	keys.ReadRSAKeys()

	app := fiber.New(fiber.Config{
		AppName: "Auth Service",
	})

	db.ConncetDB()

	go consumers.ConsumerAnswerTokens()

	// consumers.AnswerTokens()
	handler.Handlers(app)

	app.Listen(":" + config.AppConfig.Server.Port)
}
