package consumers

import (
	"coreservice/internal/di"
	"coreservice/internal/entity"
	"encoding/json"
	"fmt"

	dailyTaskRepo "coreservice/internal/adapter/db/postgres/daily_task"
)

type DailyTaskConsumer struct {
	consumer di.KafkaConsumer
	logger   di.LoggerType
	cfg      di.ConfigType
	repo     *dailyTaskRepo.DailyTaskRepo
}

func New(cfg di.ConfigType, logger di.LoggerType, consumer di.KafkaConsumer, repo *dailyTaskRepo.DailyTaskRepo) *DailyTaskConsumer {
	return &DailyTaskConsumer{
		consumer: consumer,
		logger:   logger,
		cfg:      cfg,
		repo:     repo,
	}
}

func (c *DailyTaskConsumer) Close() {
	c.consumer.Close()
}

func (c *DailyTaskConsumer) ReceiveDailyTask() {

	var err error
	var dailyTask entity.DailyTask

	c.logger.Info("RecieveSeasonInfo is working")
	err = c.consumer.Subscribe(c.cfg.Kafka.DailyTaskTopicRecieve, nil)
	if err != nil {
		c.logger.Fatal("Failed to assign partition:" + err.Error())
		panic(err)
	}

	for {
		msg, err := c.consumer.ReadMessage(-1)

		if err == nil {

			c.logger.Info("Received message: " + string(msg.Value) + " from topic:" + msg.TopicPartition.String() + " with offset " + msg.TopicPartition.Offset.String())

			err := json.Unmarshal(msg.Value, &dailyTask)
			c.logger.Info("Request recieved: " + fmt.Sprintln(dailyTask))
			if err != nil {
				c.logger.Error("Error while consuming: " + err.Error())
			}

			c.logger.Info(fmt.Sprint(dailyTask))

			c.repo.AddDailyTask(dailyTask)

		} else {
			c.logger.Error("Error while consuming: " + err.Error())
		}
	}

}
