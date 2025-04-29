package producers

import (
	"adminservice/internal/di"
	"adminservice/internal/entity"
	utils "adminservice/pkg/kafka_utils"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func SendSeasonInfo(season entity.SeasonOut, cfg di.ConfigType, bus di.Bus) {
	log.Println("Season Producer created successfully")

	message := utils.CreateMessage(season, cfg.Kafka.SeasonTopicSend, cfg.Kafka.Partition)

	deliveryChan := make(chan kafka.Event)

	err := bus.SeasonProducer.Produce(&message, deliveryChan)
	if err != nil {
		log.Printf("Failed to produce message: %s", err)
		return
	}

	log.Println("Message sent, waiting for delivery confirmation...")

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

	bus.SeasonProducer.Flush(1000)
}
