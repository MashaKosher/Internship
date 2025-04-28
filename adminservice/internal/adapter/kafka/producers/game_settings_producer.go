package producers

import (
	"adminservice/internal/config"
	"adminservice/internal/entity"
	"adminservice/pkg"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func SendGameSettings(season entity.SettingsJson) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer p.Close()

	log.Println("Season Producer created successfully")

	message := pkg.CreateMessage(season, config.AppConfig.Kafka.GameSettingsTopicSend, -1)

	deliveryChan := make(chan kafka.Event)

	err = p.Produce(&message, deliveryChan)
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

	p.Flush(1000)
}

// func SendGameSettings(season entity.SettingsJson) {
// 	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
// 	if err != nil {
// 		log.Fatalf("Failed to create producer: %s", err)
// 	}
// 	defer p.Close()

// 	log.Println("Season Producer created successfully")

// 	message := pkg.CreateMessage(season, config.AppConfig.Kafka.SeasonTopicSend, config.AppConfig.Kafka.Partition)

// 	deliveryChan := make(chan kafka.Event)

// 	err = p.Produce(&message, deliveryChan)
// 	if err != nil {
// 		log.Printf("Failed to produce message: %s", err)
// 		return
// 	}

// 	log.Println("Message sent, waiting for delivery confirmation...")

// 	go func() {
// 		event := <-deliveryChan
// 		switch e := event.(type) {
// 		case *kafka.Message:
// 			if e.TopicPartition.Error != nil {
// 				log.Printf("Delivery failed: %v", e.TopicPartition.Error)
// 			} else {
// 				log.Printf("Delivered message to %v [%d] at offset %v",
// 					*e.TopicPartition.Topic, e.TopicPartition.Partition, e.TopicPartition.Offset)
// 			}
// 		default:
// 			log.Printf("Ignored event: %s", e)
// 		}
// 	}()

// 	p.Flush(1000)
// }
