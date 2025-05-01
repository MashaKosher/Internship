package producers

import (
	"adminservice/internal/di"
	"adminservice/internal/entity"
	utils "adminservice/pkg/kafka_utils"
	"fmt"

	// "log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func CheckToken(accessToken, refreshToken string, cfg di.ConfigType, bus di.Bus) {
	bus.Logger.Info("Producer created successfully")

	var request entity.AuthRequest
	request.Partition = cfg.Kafka.Partition
	request.AccessToken = accessToken
	request.RefreshToken = refreshToken

	message := utils.CreateMessage(request, cfg.Kafka.AuthTopicRecieve, cfg.Kafka.Partition)

	// Канал для получения событий доставки
	deliveryChan := make(chan kafka.Event)

	err := bus.AuthProducer.Produce(&message, deliveryChan)
	if err != nil {
		bus.Logger.Error("Failed to produce message: " + err.Error())
		return
	}

	bus.Logger.Info("Message sent, waiting for delivery confirmation...")

	// Ожидаем подтверждения доставки
	go func() {
		event := <-deliveryChan
		switch e := event.(type) {
		case *kafka.Message:
			if e.TopicPartition.Error != nil {
				bus.Logger.Error("Delivery failed: " + e.TopicPartition.Error.Error())
			} else {
				bus.Logger.Info(fmt.Sprintf("Delivered message to %v [%d] at offset %v",
					*e.TopicPartition.Topic, e.TopicPartition.Partition, e.TopicPartition.Offset))
			}
		default:
			bus.Logger.Error("Ignored event: " + e.String())
		}
	}()

	bus.AuthProducer.Flush(1000)
}
