package consumers

import (
	"authservice/internal/config"
	"authservice/internal/entity"
	"authservice/internal/logger"
	"authservice/pkg"
	"authservice/pkg/brokercheck"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ConsumerAnswerTokens() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.AppConfig.Kafka.Host + ":" + config.AppConfig.Kafka.Port,
		"group.id":          "authService",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		logger.Logger.Fatal("Failed to create consumer: " + err.Error())
	}
	defer c.Close()

	logger.Logger.Info("Kafka Answer Auth Tokens Consumer connected successfully")

	if err := c.SubscribeTopics([]string{config.AppConfig.Kafka.TopicRecieve}, nil); err != nil {
		logger.Logger.Fatal("Failed to subscribe to topics: " + err.Error())
	}

	var authRequest entity.AuthRequest

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {

			// Тут поменять и вообще все кафка утилс
			authRequest, err := pkg.DeserializeAuthAnswer(msg.Value, authRequest)

			if err != nil {
				logger.Logger.Fatal("Error while consuming: " + err.Error())
			}

			logger.Logger.Info("Received message: " + string(msg.Value) + " from topic:" + fmt.Sprintln(msg.TopicPartition))
			go brokercheck.BrokerCheck(authRequest)
		} else {
			logger.Logger.Fatal("Error while consuming: " + err.Error())
		}
	}
}
