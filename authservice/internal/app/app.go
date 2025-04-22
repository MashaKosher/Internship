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

	// httpServer := httpserver.New(httpserver.Port(config.AppConfig.Server.Host))

	// Получаем сущность DB
	db := db.ConncetDB()

	// Создаем Use Case
	authUseCase := auth.New(
		authRepo.New(db), // создаем конкретный репозиторий и передаем в конретный Use Case
	)

	go consumers.ConsumerAnswerTokens()

	v1.NewRouter(app, authUseCase)

	// // Waiting signal
	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// select {
	// case s := <-interrupt:
	// 	logger.Logger.Info("app - Run - signal: " + s.String())
	// case err := <-httpServer.Notify():
	// 	logger.Logger.Error("app - Run - httpServer.Notify: " + err.Error())

	// }

	// // Shutdown
	// err := httpServer.Shutdown()
	// if err != nil {
	// 	logger.Logger.Error("app - Run - httpServer.Shutdown: " + err.Error())
	// }

	// app.Listen(":" + config.AppConfig.Server.Port)

	go func() {
		if err := app.Listen(":" + config.AppConfig.Server.Port); err != nil {
			logger.Logger.Error("Server error:" + err.Error())
		}
	}()

	// Канал для получения сигналов ОС
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // Перехватываем SIGINT и SIGTERM
	<-quit                                               // Ждем сигнал завершения

	logger.Logger.Info("Shutting down server...")

	// Создаем контекст с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Останавливаем сервер
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Logger.Fatal("Server forced to shutdown: " + err.Error())
	}

	logger.Logger.Info("Server gracefully stopped")
}
