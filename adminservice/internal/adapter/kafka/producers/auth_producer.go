package producers

import (
	"adminservice/internal/di"
	"adminservice/internal/entity"
	"adminservice/pkg"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func CheckToken(accessToken, refreshToken string, cfg di.ConfigType, bus di.Bus) {
	// p, err := kafka.NewProducer(&kafka.ConfigMap{
	// 	"bootstrap.servers": config.AppConfig.Kafka.Host + ":" + config.AppConfig.Kafka.Port, // Используйте localhost
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to create producer: %s", err)
	// }
	// defer p.Close()

	log.Println("Producer created successfully")

	var request entity.AuthRequest
	request.Partition = cfg.Kafka.Partition
	request.AccessToken = accessToken
	request.RefreshToken = refreshToken

	message := pkg.CreateMessage(request, cfg.Kafka.AuthTopicRecieve, cfg.Kafka.Partition)

	// Канал для получения событий доставки
	deliveryChan := make(chan kafka.Event)

	err := bus.AuthProducer.Produce(&message, deliveryChan)
	if err != nil {
		log.Printf("Failed to produce message: %s", err)
		return
	}

	log.Println("Message sent, waiting for delivery confirmation...")

	// Ожидаем подтверждения доставки
	go func() {
		event := <-deliveryChan
		switch e := event.(type) {
		case *kafka.Message:
			if e.TopicPartition.Error != nil {
				log.Printf("Delivery failed: %v", e.TopicPartition.Error)
			} else {
				log.Printf("Delivered message to %v [%d] at offset %v",
					*e.TopicPartition.Topic, e.TopicPartition.Partition, e.TopicPartition.Offset)
			}
		default:
			log.Printf("Ignored event: %s", e)
		}
	}()

	bus.AuthProducer.Flush(1000)
}
