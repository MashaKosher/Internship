package producers

import (
	"fmt"
	"gameservice/internal/config"
	"gameservice/internal/entity"
	"gameservice/pkg"
	"gameservice/pkg/logger"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func CheckToken(accessToken, refreshToken string) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": config.AppConfig.Kafka.Host + ":" + config.AppConfig.Kafka.Port,
	})
	if err != nil {
		logger.L.Fatal("Failed to create producer: " + err.Error())
	}
	defer p.Close()

	logger.L.Info("Producer created successfully")

	var request entity.AuthRequest
	request.Partition = config.AppConfig.Kafka.Partition
	request.AccessToken = accessToken
	request.RefreshToken = refreshToken

	message := pkg.CreateMessage(request, config.AppConfig.Kafka.AuthTopicRecieve, config.AppConfig.Kafka.Partition)

	// Канал для получения событий доставки
	deliveryChan := make(chan kafka.Event)

	err = p.Produce(&message, deliveryChan)
	if err != nil {
		logger.L.Error("Failed to produce message: " + err.Error())
		return
	}

	logger.L.Info("Message sent, waiting for delivery confirmation...")

	// Ожидаем подтверждения доставки
	go func() {
		event := <-deliveryChan
		switch e := event.(type) {
		case *kafka.Message:
			if e.TopicPartition.Error != nil {
				logger.L.Error("Delivery failed: " + e.TopicPartition.Error.Error())
			} else {
				logger.L.Info(fmt.Sprintf("Delivered message to %v [%d] at offset %v",
					*e.TopicPartition.Topic, e.TopicPartition.Partition, e.TopicPartition.Offset))
			}
		default:
			logger.L.Info("Ignored event: " + e.String())
		}
	}()

	p.Flush(1000)
}
