package consumers

import (
	"coreservice/internal/config"
	"coreservice/internal/entity"
	"coreservice/internal/logger"
	"encoding/json"
	"fmt"
	"log"

	repo "coreservice/internal/repository/sqlc"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ReceiveDailyTask() {

	var err error
	var dailyTask entity.DailyTask

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.AppConfig.Kafka.Host + ":" + config.AppConfig.Kafka.Port, // Используйте localhost
		"group.id":          "authService",
		"auto.offset.reset": "latest",
	})
	if err != nil {
		logger.Logger.Error("Failed to create consumer: " + err.Error())
	}
	defer consumer.Close()

	logger.Logger.Info("RecieveSeasonInfo is working")
	err = consumer.Subscribe(config.AppConfig.Kafka.DailyTaskTopicRecieve, nil)
	if err != nil {
		logger.Logger.Fatal("Failed to assign partition:" + err.Error())
		panic(err)
	}

	for {
		msg, err := consumer.ReadMessage(-1)

		if err == nil {

			log.Println("Maessage readed")
			logger.Logger.Info("Received message: " + string(msg.Value) + " from topic:" + msg.TopicPartition.String() + " with offset " + msg.TopicPartition.Offset.String())

			err := json.Unmarshal(msg.Value, &dailyTask)
			logger.Logger.Info("Request recieved: " + fmt.Sprintln(dailyTask))
			if err != nil {
				logger.Logger.Error("Error while consuming: " + err.Error())
			}

			logger.Logger.Info(fmt.Sprint(dailyTask))

			repo.AddDailyTask(dailyTask)

		} else {
			logger.Logger.Error("Error while consuming: " + err.Error())
		}
	}

}
