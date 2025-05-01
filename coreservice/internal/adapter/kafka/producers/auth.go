package producers

import (
	"coreservice/internal/di"
	"coreservice/internal/entity"
	"coreservice/pkg"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func CheckToken(accessToken, resfreshToken string, cfg di.ConfigType, bus di.Bus) {
	// producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": config.AppConfig.Kafka.Host + ":" + config.AppConfig.Kafka.Port})
	// if err != nil {
	// 	logger.Logger.Fatal("Failed to create producer:" + err.Error())
	// }
	// defer producer.Close()

	bus.Logger.Info("Auth Producer created successfully")

	var request entity.AuthRequest
	request.Partition = cfg.Kafka.Partition
	request.AccessToken = accessToken
	request.RefreshToken = resfreshToken

	message := pkg.CreateMessage(request, cfg.Kafka.AuthTopicSend, cfg.Kafka.Partition)

	//	Канал для получения событий доставки
	deliveryChan := make(chan kafka.Event)

	err := bus.AuthProducer.Produce(&message, deliveryChan)
	if err != nil {
		bus.Logger.Error("Failed to produce message:" + err.Error())
		return
	}

	bus.Logger.Info("Message sent, waiting for delivery confirmation...")

	go func() {
		event := <-deliveryChan
		switch e := event.(type) {
		case *kafka.Message:
			if e.TopicPartition.Error != nil {
				bus.Logger.Fatal("Delivery failed: " + fmt.Sprintln(e.TopicPartition.Error))
			} else {
				bus.Logger.Info(fmt.Sprintf("Delivered message to %v [%d] at offset %v",
					*e.TopicPartition.Topic, e.TopicPartition.Partition, e.TopicPartition.Offset))
			}
		default:
			bus.Logger.Info("Ignored event: " + fmt.Sprintln(e))
		}
	}()

	bus.AuthProducer.Flush(1000)
}
