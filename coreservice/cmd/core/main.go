package main

import (
	"coreservice/internal/adapter/asynq/consumer"
	"coreservice/internal/adapter/elastic"
	"coreservice/internal/adapter/kafka/consumers"
	"coreservice/internal/config"
	v1 "coreservice/internal/controller/http"
	"coreservice/internal/di/setup"
	"coreservice/internal/logger"
	"coreservice/internal/repository/sqlc"

	"github.com/gin-gonic/gin"

	_ "coreservice/docs"
)

// @title Example API
// @version 1.0
// @description This is a sample API for demonstrating Swagger with Gin.
// @host localhost:8006
// @BasePath /
func main() {
	cfg := config.MustParseConfig()
	deps := setup.MustContainer(cfg)

	logFile := logger.CreateLogger()
	defer logFile.Close()
	defer logger.Logger.Sync()
	defer logger.Logger.Info("Program end")

	conn, ctx := sqlc.DBConnect()
	defer conn.Close(ctx)

	go consumer.AsynqConsumer()

	go consumers.RecieveSeasonInfo()
	go consumers.ReceiveDailyTask()

	elastic.ESClientConnection()
	elastic.ESCreateIndexIfNotexist()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	v1.NewRouter(router, deps)

	router.Run(config.AppConfig.Server.Host + ":" + config.AppConfig.Server.Port)
}
