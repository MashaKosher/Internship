package app

import (
	"coreservice/internal/adapter/asynq/consumer"
	"coreservice/internal/adapter/kafka/consumers"
	"coreservice/internal/di"
	"coreservice/internal/di/setup"

	v1 "coreservice/internal/controller/http"

	"github.com/gin-gonic/gin"
)

func Run(cfg di.ConfigType) {

	deps := setup.MustContainer(cfg)
	go consumer.AsynqConsumer()

	go consumers.RecieveSeasonInfo(cfg, deps.Bus, deps.DB, deps.Elastic.ESClient, deps.Elastic.SeasonSearchIndex)
	go consumers.ReceiveDailyTask(cfg, deps.Bus)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	v1.NewRouter(router, deps)
	router.Run(cfg.Server.Host + ":" + cfg.Server.Port)
}
