package consumers

import (
	"coreservice/internal/di"
	"coreservice/internal/entity"
	repo "coreservice/internal/repository/sqlc"
	"encoding/json"
	"fmt"
)

func ReceiveDailyTask(cfg di.ConfigType, bus di.Bus) {

	var err error
	var dailyTask entity.DailyTask

	// consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
	// 	"bootstrap.servers": config.AppConfig.Kafka.Host + ":" + config.AppConfig.Kafka.Port, // Используйте localhost
	// 	"group.id":          "authService",
	// 	"auto.offset.reset": "latest",
	// })
	// if err != nil {
	// 	logger.Logger.Error("Failed to create consumer: " + err.Error())
	// }
	// defer consumer.Close()

	bus.Logger.Info("RecieveSeasonInfo is working")
	err = bus.DailyTaskConsumer.Subscribe(cfg.Kafka.DailyTaskTopicRecieve, nil)
	if err != nil {
		bus.Logger.Fatal("Failed to assign partition:" + err.Error())
		panic(err)
	}

	for {
		msg, err := bus.DailyTaskConsumer.ReadMessage(-1)

		if err == nil {

			bus.Logger.Info("Received message: " + string(msg.Value) + " from topic:" + msg.TopicPartition.String() + " with offset " + msg.TopicPartition.Offset.String())

			err := json.Unmarshal(msg.Value, &dailyTask)
			bus.Logger.Info("Request recieved: " + fmt.Sprintln(dailyTask))
			if err != nil {
				bus.Logger.Error("Error while consuming: " + err.Error())
			}

			bus.Logger.Info(fmt.Sprint(dailyTask))

			repo.AddDailyTask(dailyTask)

		} else {
			bus.Logger.Error("Error while consuming: " + err.Error())
		}
	}

}
