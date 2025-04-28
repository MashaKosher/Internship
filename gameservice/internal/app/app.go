package app

import (
	"context"
	gameRepo "gameservice/internal/adapter/db/clickhouse/game"
	"gameservice/internal/adapter/kafka/consumers"
	redisRepo "gameservice/internal/adapter/redis/game_settings"
	"gameservice/internal/config"
	v1 "gameservice/internal/controller"
	"gameservice/pkg/client/clickhouse"
	"gameservice/pkg/client/redis"
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

	redisDB, ctx := redis.InitRedis()

	// Здесь нужна БД
	db := clickhouse.ConnectClickHouse()

	gameRepository := gameRepo.New(db)
	redisRepository := redisRepo.New(redisDB, ctx)

	if err := gameRepository.CreateGameResultsTable(); err != nil {
		logger.L.Fatal("Problems with creating tables")
	}

	logger.L.Info("ClickHouse table created successfully")

	gameUseCase := game.New(
		gameRepository,
		redisRepository,
	)

	go consumers.GameSettingsConsumer(redisRepository)

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
