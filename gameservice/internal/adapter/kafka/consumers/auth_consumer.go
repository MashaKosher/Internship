package consumers

import (
	"gameservice/internal/config"
	"gameservice/internal/entity"
	"gameservice/pkg"
	"gameservice/pkg/logger"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func RecieveTokenInfo() (entity.AuthAnswer, error) {

	var err error
	var answer entity.AuthAnswer

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.AppConfig.Kafka.Host + ":" + config.AppConfig.Kafka.Port, // Используйте localhost
		"group.id":          "adminServiceBecomingAnswer",
		"auto.offset.reset": "latest",
	})
	if err != nil {
		logger.L.Fatal("Failed to create consumer:" + err.Error())
	}
	defer consumer.Close()
	logger.L.Info("Kafka connected successfully")

	err = consumer.Assign([]kafka.TopicPartition{{Topic: &config.AppConfig.Kafka.AuthTopicSend, Partition: config.AppConfig.Kafka.Partition, Offset: kafka.OffsetTail(1)}})
	if err != nil {
		logger.L.Fatal("Failed to assign partition:" + err.Error())
	}

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			logger.L.Info("Received message: " + string(msg.Value) + " from topic:" + msg.TopicPartition.String() + " with offset " + msg.TopicPartition.Offset.String())

			answer, err := pkg.DeserializeAuthAnswer(msg.Value, answer)
			if err != nil {
				logger.L.Error("Error while consuming:" + err.Error())
				return answer, err
			}
			return answer, err
		} else {
			logger.L.Error("Error while consuming: " + err.Error())
			return answer, err
		}
	}
}
