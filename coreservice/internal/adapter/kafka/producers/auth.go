package producers

import (
	// "adminservice/internal/entity"

	"coreservice/internal/config"
	"coreservice/internal/entity"
	"coreservice/internal/logger"
	"coreservice/pkg"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func CheckToken(accessToken, resfreshToken string) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": config.AppConfig.Kafka.Host + ":" + config.AppConfig.Kafka.Port})
	if err != nil {
		logger.Logger.Fatal("Failed to create producer:" + err.Error())
	}
	defer producer.Close()

	logger.Logger.Info("Producer created successfully")

	var request entity.AuthRequest
	request.Partition = config.AppConfig.Kafka.Partition
	request.AccessToken = accessToken
	request.RefreshToken = resfreshToken

	message := pkg.CreateMessage(request, config.AppConfig.Kafka.AuthTopicSend, config.AppConfig.Kafka.Partition)

	//	Канал для получения событий доставки
	deliveryChan := make(chan kafka.Event)

	err = producer.Produce(&message, deliveryChan)
	if err != nil {
		logger.Logger.Error("Failed to produce message:" + err.Error())
		return
	}

	logger.Logger.Info("Message sent, waiting for delivery confirmation...")

	go func() {
		event := <-deliveryChan
		switch e := event.(type) {
		case *kafka.Message:
			if e.TopicPartition.Error != nil {
				logger.Logger.Fatal("Delivery failed: " + fmt.Sprintln(e.TopicPartition.Error))
			} else {
				logger.Logger.Info(fmt.Sprintf("Delivered message to %v [%d] at offset %v",
					*e.TopicPartition.Topic, e.TopicPartition.Partition, e.TopicPartition.Offset))
			}
		default:
			logger.Logger.Info("Ignored event: " + fmt.Sprintln(e))
		}
	}()

	producer.Flush(1000)
}
