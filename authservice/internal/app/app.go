package app

import (
	"authservice/internal/adapter/kafka/consumers"
	"authservice/internal/config"
	v1 "authservice/internal/controller/http"
	db "authservice/pkg/client/sql"
	"authservice/pkg/keys"
	"authservice/pkg/logger"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	go consumers.ConsumerAnswerTokens(authUseCase)

	v1.NewRouter(app, authUseCase)

	go func() {
		if err := app.Listen(":" + config.AppConfig.Server.Port); err != nil {
			logger.Logger.Error("Server error:" + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Logger.Fatal("Server forced to shutdown: " + err.Error())
	}

	logger.Logger.Info("Server gracefully stopped")
}
