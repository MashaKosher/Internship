package consumer

import (
	"coreservice/internal/adapter/asynq/tasks"
	"coreservice/internal/config"
	"coreservice/internal/logger"
	"strings"
	"time"

	"github.com/hibiken/asynq"
)

func AsynqConsumer() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     config.AppConfig.Redis.Host + ":" + config.AppConfig.Redis.Port,
			PoolSize: 10,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 1,
			},
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				return 2 * time.Second
			},
			IsFailure: func(err error) bool {
				if strings.Contains(err.Error(), "task expired") {
					return false
				}
				return true
			},
		},
	)

	logger.Logger.Info("Asynq consumer connected successfully")
	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeSeason, tasks.SeasonTaskHadler)

	if err := srv.Run(mux); err != nil {
		logger.Logger.Fatal("could not run server: " + err.Error())
	}

	logger.Logger.Info("Asynq consumer dont work")
}
