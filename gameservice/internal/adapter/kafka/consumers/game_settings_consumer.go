package consumers

import (
	"fmt"
	"gameservice/internal/config"
	"gameservice/pkg"
	"gameservice/pkg/logger"

	redisRepo "gameservice/internal/adapter/redis/game_settings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func GameSettingsConsumer(redisRepo *redisRepo.GameSettingsRepo) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.AppConfig.Kafka.Host + ":" + config.AppConfig.Kafka.Port,
		"group.id":          "authService",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		logger.L.Fatal("Failed to create consumer: " + err.Error())
	}
	defer c.Close()

	logger.L.Info("Kafka Game Settings Consumer connected successfully")

	if err := c.SubscribeTopics([]string{config.AppConfig.Kafka.GameSettingsTopicRecieve}, nil); err != nil {
		logger.L.Fatal("Failed to subscribe to topics: " + err.Error())
	}

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {

			// Тут поменять и вообще все кафка утилс
			gameSettings, err := pkg.DeserializeGameSettings(msg.Value)

			if err != nil {
				logger.L.Fatal("Error while consuming: " + err.Error())
			}

			logger.L.Info("Received message: " + string(msg.Value) + " from topic:" + fmt.Sprintln(msg.TopicPartition))

			if err := redisRepo.RefreshGameSettings(gameSettings); err != nil {
				logger.L.Fatal("Troubles with adding game setiings to Redis: " + err.Error())
			}

			logger.L.Info("Game Settings added successfuly: " + fmt.Sprint(gameSettings))

		} else {
			logger.L.Fatal("Error while consuming: " + err.Error())
		}
	}
}
