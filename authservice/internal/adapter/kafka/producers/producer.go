package producers

import (
	"authservice/internal/config"
	"authservice/internal/entity"
	"authservice/internal/logger"
	"authservice/pkg"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func AnswerToken(answer entity.AuthAnswer, partition int32) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": config.AppConfig.Kafka.Host + ":" + config.AppConfig.Kafka.Port})
	if err != nil {
		logger.Logger.Error("Failed to create producer:" + err.Error())
	}
	defer p.Close()

	logger.Logger.Info("Producer created successfully")

	message := pkg.CreateMessage(answer, config.AppConfig.Kafka.TopicSend, partition)

	deliveryChan := make(chan kafka.Event)

	err = p.Produce(&message, deliveryChan)
	if err != nil {
		logger.Logger.Error("Failed to produce message: " + err.Error())
		panic(err)

	}

	go func() {
		event := <-deliveryChan
		switch e := event.(type) {
		case *kafka.Message:
			if e.TopicPartition.Error != nil {
				logger.Logger.Error(fmt.Sprintf("Delivery failed: %v", e.TopicPartition.Error))
			} else {
				logger.Logger.Error(fmt.Sprintf("Delivered message to %v [%d] at offset %v",
					*e.TopicPartition.Topic, e.TopicPartition.Partition, e.TopicPartition.Offset))
			}
		default:
			logger.Logger.Error(fmt.Sprintf("Ignored event: %s", e))
		}
	}()

	p.Flush(1000)
}
