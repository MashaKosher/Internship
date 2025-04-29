package producers

import (
	"authservice/internal/di"
	"authservice/internal/entity"
	"authservice/pkg"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func AnswerToken(answer entity.AuthAnswer, partition int32, logger di.LoggerType, cfg di.ConfigType, bus di.Bus) {
	logger.Info("Producer created successfully")

	message := pkg.CreateMessage(answer, cfg.Kafka.TopicSend, partition, logger)

	deliveryChan := make(chan kafka.Event)

	err := bus.Producer.Produce(&message, deliveryChan)
	if err != nil {
		logger.Fatal("Failed to produce message: " + err.Error())
	}

	go func() {
		event := <-deliveryChan
		switch e := event.(type) {
		case *kafka.Message:
			if e.TopicPartition.Error != nil {
				logger.Error(fmt.Sprintf("Delivery failed: %v", e.TopicPartition.Error))
			} else {
				logger.Error(fmt.Sprintf("Delivered message to %v [%d] at offset %v",
					*e.TopicPartition.Topic, e.TopicPartition.Partition, e.TopicPartition.Offset))
			}
		default:
			logger.Error(fmt.Sprintf("Ignored event: %s", e))
		}
	}()

	bus.Producer.Flush(1000)
}
