package consumer

import (
	"coreservice/internal/adapter/asynq/tasks"
	"coreservice/internal/config"
	"coreservice/internal/di"
	"strings"
	"time"

	"github.com/hibiken/asynq"
)

const connPoolsize = 10
const concSize = 10
const criticalPriority = 1
const retryTimeDelay = 2 * time.Second

func AsynqConsumer(deps di.Container) {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     config.AppConfig.Redis.Host + ":" + config.AppConfig.Redis.Port,
			PoolSize: connPoolsize,
		},
		asynq.Config{
			Concurrency: concSize,
			Queues: map[string]int{
				"critical": criticalPriority,
			},
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				return retryTimeDelay
			},
			IsFailure: func(err error) bool {
				if strings.Contains(err.Error(), "task expired") {
					return false
				}
				return true
			},
		},
	)

	deps.Logger.Info("Asynq consumer connected successfully")
	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeSeason, tasks.SeasonTaskHadler(deps.Logger, deps.DB, deps.Elastic.ESClient, deps.Elastic.SeasonSearchIndex))

	if err := srv.Run(mux); err != nil {
		deps.Logger.Fatal("could not run server: " + err.Error())
	}

	deps.Logger.Info("Asynq consumer dont work")
}
