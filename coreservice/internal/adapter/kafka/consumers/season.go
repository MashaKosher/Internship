package consumers

import (
	"coreservice/internal/di"
	"coreservice/internal/entity"
	"coreservice/pkg"
	"fmt"

	"coreservice/internal/adapter/asynq/producer"
	seasonRepo "coreservice/internal/adapter/db/postgres/season"
	seasonStatusElasticRepo "coreservice/internal/adapter/elastic/seasons"
)

func RecieveSeasonInfo(cfg di.ConfigType, bus di.Bus, db di.DBType, ESClient di.ESClient, Index di.ElasticIndex) {

	seasonRepo := seasonRepo.New(db)
	elastic := seasonStatusElasticRepo.New(ESClient, Index, bus.Logger)

	var err error
	var season entity.Season

	// consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
	// 	"bootstrap.servers": cfg.Kafka.Host + ":" + cfg.Kafka.Port, // Используйте localhost
	// 	"group.id":          "authService",
	// 	"auto.offset.reset": "latest",
	// })
	// if err != nil {
	// 	logger.Logger.Error("Failed to create consumer: " + err.Error())
	// }
	// defer consumer.Close()

	bus.Logger.Info("RecieveSeasonInfo is working")
	err = bus.SeasonConsumer.Subscribe(cfg.Kafka.SeasonTopicRecieve, nil)
	if err != nil {
		bus.Logger.Fatal("Failed to assign partition:" + err.Error())
		panic(err)
	}

	for {
		msg, err := bus.SeasonConsumer.ReadMessage(-1)

		if err == nil {

			bus.Logger.Info("Received message: " + string(msg.Value) + " from topic:" + msg.TopicPartition.String() + " with offset " + msg.TopicPartition.Offset.String())

			answer, err := pkg.DeseriSeasonAnswer(msg.Value, season)
			if err != nil {
				bus.Logger.Error("Error while consuming: " + err.Error())
			}

			bus.Logger.Info(fmt.Sprint(answer))

			// add to db
			if err = seasonRepo.AddSeason(answer); err != nil {
				panic(err)
			}

			bus.Logger.Info("Start date: " + answer.StartDate + " End date: " + answer.EndDate)

			startTime, err := pkg.ParseTimeToLocal(answer.StartDate)
			if err != nil {
				panic(err)
			}

			endTime, err := pkg.ParseTimeToLocal(answer.EndDate)
			if err != nil {
				panic(err)
			}
			producer.PlanSeasonTasks(int(answer.ID), startTime, endTime)

			elastic.AddSeasonToIndex(int(answer.ID))
		} else {
			bus.Logger.Error("Error while consuming: " + err.Error())
		}
	}
}
