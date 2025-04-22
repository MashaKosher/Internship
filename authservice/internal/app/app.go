package app

import (
	"authservice/internal/adapter/kafka/consumers"
	"authservice/internal/config"
	v1 "authservice/internal/controller/http"
	db "authservice/pkg/client/sql"
	"authservice/pkg/keys"
	"authservice/pkg/logger"

	authRepo "authservice/internal/adapter/db/sql/auth"
	"authservice/internal/usecase/auth"

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

	// Получаем сущность DB
	db := db.ConncetDB()

	// Создаем Use Case
	authUseCase := auth.New(
		authRepo.New(db), // создаем конкретный репозиторий и передаем в конретный Use Case
	)

	go consumers.ConsumerAnswerTokens()

	// handler.Handlers(app)
	v1.NewRouter(app, authUseCase)

	app.Listen(":" + config.AppConfig.Server.Port)
}
