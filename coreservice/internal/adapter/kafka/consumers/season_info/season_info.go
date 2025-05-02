package consumers

import (
	"coreservice/internal/di"
	"coreservice/internal/entity"
	"coreservice/pkg"
	utils "coreservice/pkg/kafka_utils"
	"fmt"

	"coreservice/internal/adapter/asynq/producer"
	seasonRepo "coreservice/internal/adapter/db/postgres/season"
	seasonStatusElasticRepo "coreservice/internal/adapter/elastic/seasons"
)

type SeasonInfoConsumer struct {
	consumer di.KafkaConsumer
	logger   di.LoggerType
	cfg      di.ConfigType
	repo     *seasonRepo.SeasonRepo
	elastic  *seasonStatusElasticRepo.SeasonStatusRepo
}

func New(
	cfg di.ConfigType,
	logger di.LoggerType,
	consumer di.KafkaConsumer,
	repo *seasonRepo.SeasonRepo,
	elastic *seasonStatusElasticRepo.SeasonStatusRepo,
) *SeasonInfoConsumer {
	return &SeasonInfoConsumer{
		consumer: consumer,
		logger:   logger,
		cfg:      cfg,
		repo:     repo,
		elastic:  elastic,
	}
}

func (c *SeasonInfoConsumer) Close() {
	c.consumer.Close()
}

func (c *SeasonInfoConsumer) RecieveSeasonInfo() {
	var err error
	var season entity.Season

	c.logger.Info("RecieveSeasonInfo is working")
	err = c.consumer.Subscribe(c.cfg.Kafka.SeasonTopicRecieve, nil)
	if err != nil {
		c.logger.Fatal("Failed to assign partition:" + err.Error())
		panic(err)
	}

	for {
		msg, err := c.consumer.ReadMessage(-1)

		if err == nil {

			c.logger.Info("Received message: " + string(msg.Value) + " from topic:" + msg.TopicPartition.String() + " with offset " + msg.TopicPartition.Offset.String())

			answer, err := utils.DeseriSeasonAnswer(msg.Value, season, c.logger)
			if err != nil {
				c.logger.Error("Error while consuming: " + err.Error())
			}

			c.logger.Info(fmt.Sprint(answer))

			// add to db
			if err = c.repo.AddSeason(answer); err != nil {
				panic(err)
			}

			c.logger.Info("Start date: " + answer.StartDate + " End date: " + answer.EndDate)

			startTime, err := pkg.ParseTimeToLocal(answer.StartDate)
			if err != nil {
				panic(err)
			}

			endTime, err := pkg.ParseTimeToLocal(answer.EndDate)
			if err != nil {
				panic(err)
			}
			producer.PlanSeasonTasks(int(answer.ID), startTime, endTime, c.cfg, c.logger)

			c.elastic.AddSeasonToIndex(int(answer.ID))
		} else {
			c.logger.Error("Error while consuming: " + err.Error())
		}
	}
}
