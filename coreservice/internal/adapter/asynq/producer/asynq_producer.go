package producer

import (
	"coreservice/internal/adapter/asynq/tasks"
	"coreservice/internal/di"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
)

const clientPoolSize = 10
const maxSendRetry = 3
const queueType = "critical"

func PlanSeasonTasks(seasonID int, startTime, endTime time.Time, cfg di.ConfigType, logger di.LoggerType) {

	client := asynq.NewClient(
		asynq.RedisClientOpt{
			Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
			PoolSize: clientPoolSize,
		})
	defer client.Close()

	logger.Info("Asynq Client  connected successfully")

	// Creating task for season start
	task, err := tasks.NewSeasonTask(seasonID, startTime, tasks.CurrentSeason)

	if err != nil {
		panic(err)
	}
	logger.Info("Start Time dif: " + fmt.Sprint(time.Until(startTime)))

	info, err := client.Enqueue(
		task,
		asynq.ProcessAt(startTime),
		asynq.MaxRetry(maxSendRetry),
		asynq.Queue(queueType),
	)

	logger.Info("Task will be start at: " + fmt.Sprint(startTime))
	if err != nil {
		panic(err)
	}

	logger.Info(fmt.Sprintf("enqued task: id = %s queue = %s ", info.ID, info.Queue))

	// creating task for end season
	task, err = tasks.NewSeasonTask(seasonID, endTime, tasks.CancelSeason)

	if err != nil {
		panic(err)
	}

	logger.Info("End Time dif: " + fmt.Sprint(time.Until(endTime)))

	info, err = client.Enqueue(
		task,
		// asynq.ProcessAt(time.Now().Add(time.Second*30)),
		asynq.ProcessAt(endTime),
		asynq.MaxRetry(maxSendRetry),
		asynq.Queue(queueType),
	)
	logger.Info("Task will be end at: " + fmt.Sprint(endTime))
	if err != nil {
		panic(err)
	}

	logger.Info(fmt.Sprintf("enqued task: id = %s queue = %s ", info.ID, info.Queue))
}
