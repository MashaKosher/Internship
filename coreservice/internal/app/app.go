package app

import (
	"coreservice/internal/adapter/asynq/consumer"
	"coreservice/internal/di"
	"coreservice/internal/di/setup"

	v1 "coreservice/internal/controller/http"

	"github.com/gin-gonic/gin"
)

func Run(cfg di.ConfigType) {

	deps := setup.MustContainer(cfg)
	go consumer.AsynqConsumer(deps)
	go deps.Bus.SeasonInfoConsumer.RecieveSeasonInfo()
	go deps.Bus.DailyTaskConsumer.ReceiveDailyTask()
	go deps.Bus.MatchInfoConsumer.RecieveMatchInfo()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	v1.NewRouter(router, deps)
	router.Run(cfg.Server.Host + ":" + cfg.Server.Port)
}
