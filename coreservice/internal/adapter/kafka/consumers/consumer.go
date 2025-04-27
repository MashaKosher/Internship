package consumers

import (
	"coreservice/internal/config"
	"coreservice/internal/entity"
	"coreservice/internal/logger"
	"coreservice/pkg"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func RecieveTokenInfo() (entity.AuthAnswer, error) {
	var err error
	var answer entity.AuthAnswer

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.AppConfig.Kafka.Host + ":" + config.AppConfig.Kafka.Port, // Используйте localhost
		"group.id":          "authService",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		logger.Logger.Error("Failed to create consumer: " + err.Error())
	}
	defer consumer.Close()

	logger.Logger.Info("Kafka connected successdully")
	err = consumer.Assign([]kafka.TopicPartition{{Topic: &config.AppConfig.Kafka.AuthTopicRecieve, Partition: config.AppConfig.Kafka.Partition, Offset: kafka.OffsetTail(1)}})
	if err != nil {
		logger.Logger.Fatal("Failed to assign partition:" + err.Error())
	}

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			logger.Logger.Info("Received message: " + string(msg.Value) + " from topic:" + msg.TopicPartition.String() + " with offset " + msg.TopicPartition.Offset.String())

			answer, err := pkg.DeserializeAuthAnswer(msg.Value, answer)
			if err != nil {
				logger.Logger.Error("Error while consuming: " + err.Error())
				return answer, err
			}
			return answer, err
		} else {
			logger.Logger.Error("Error while consuming: " + err.Error())
			return answer, err
		}
	}
}
