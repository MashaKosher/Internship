package producer

import (
	"coreservice/internal/adapter/asynq/tasks"
	"coreservice/internal/config"
	"coreservice/internal/logger"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

func PlanSeasonTasks(seasonID int, startTime, endTime time.Time) {
	// some sron code
	// cronTest()

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: config.AppConfig.Redis.Host + ":" + config.AppConfig.Redis.Port, PoolSize: 10})
	defer client.Close()

	logger.Logger.Info("Asynq Client  connected successfully")

	// season := ebSeason{SeasonID: rand.Int(), StartTime: time.Now().Add(time.Second * 10), EndTime: time.Now().Add(time.Second * 20)}

	// Creating task for season start
	task, err := tasks.NewSeasonTask(seasonID, startTime, tasks.CurrentSeason)

	if err != nil {
		panic(err)
	}

	fmt.Println("Start Time dif:: ", startTime.Sub(time.Now()))

	info, err := client.Enqueue(
		task,
		// asynq.ProcessAt(time.Now().Add(time.Second*20)),
		asynq.ProcessAt(startTime),
		asynq.MaxRetry(3),
		asynq.Queue("critical"),
	)

	logger.Logger.Info("Task will be start at: " + fmt.Sprint(startTime))
	if err != nil {
		panic(err)
	}

	log.Printf("enqued task: id = %s queue = %s ", info.ID, info.Queue)

	// creating task for end season
	task, err = tasks.NewSeasonTask(seasonID, endTime, tasks.CancelSeason)

	if err != nil {
		panic(err)
	}

	fmt.Println("End Time dif:: ", endTime.Sub(time.Now()))

	info, err = client.Enqueue(
		task,
		// asynq.ProcessAt(time.Now().Add(time.Second*30)),
		asynq.ProcessAt(endTime),
		asynq.MaxRetry(3),
		asynq.Queue("critical"),
	)
	logger.Logger.Info("Task will be end at: " + fmt.Sprint(endTime))
	if err != nil {
		panic(err)
	}

	log.Printf("enqued task: id = %s queue = %s ", info.ID, info.Queue)

}
