package app

import (
	"context"
	gameRepo "gameservice/internal/adapter/db/mongo/game"
	"gameservice/internal/config"
	v1 "gameservice/internal/controller"
	mongodb "gameservice/pkg/client/mongo"
	"gameservice/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gameservice/internal/usecase/game"

	"github.com/labstack/echo/v4"
)

func Run() {

	// Creating Log File
	logFile := logger.CreateLogger()
	defer logFile.Close()
	defer logger.L.Sync()
	defer logger.L.Info("Program ended")

	ctx, client, db := mongodb.MongoConnect()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			logger.L.Fatal(err.Error())
		}

		logger.L.Info("Mongo Disconnected")
	}()

	gameUseCase := game.New(
		gameRepo.New(db),
	)

	// Graceful Shutdown
	e := echo.New()
	go func() {
		if err := e.Start(":" + config.AppConfig.Server.Port); err != nil && err != http.ErrServerClosed {
			logger.L.Fatal("Server error:" + err.Error())
		}
	}()

	v1.NewRouter(e, gameUseCase)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.L.Fatal("Forced Shutdown:" + err.Error())
	}

	logger.L.Info("Server gracefully stopped")
}
