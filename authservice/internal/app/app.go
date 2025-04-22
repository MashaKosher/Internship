package app

import (
	"authservice/internal/adapter/kafka/consumers"
	"authservice/internal/config"
	v1 "authservice/internal/controller/http"
	db "authservice/pkg/client/sql"
	"authservice/pkg/keys"
	"authservice/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func Run() {
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

	// handler.Handlers(app)
	v1.NewRouter(app)

	app.Listen(":" + config.AppConfig.Server.Port)
}
