package producers

import (
	"adminservice/internal/di"
	"adminservice/internal/entity"
	"adminservice/pkg"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func SendDailyTask(task entity.DailyTasks, cfg di.ConfigType, bus di.Bus) {
	// p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	// if err != nil {
	// 	log.Fatalf("Failed to create producer: %s", err)
	// }
	// defer p.Close()

	log.Println("Task Producer created successfully")

	message := pkg.CreateMessage(task, cfg.Kafka.DailyTaskTopicSend, cfg.Kafka.Partition)

	// Канал для получения событий доставки
	deliveryChan := make(chan kafka.Event)

	err := bus.TaskProducer.Produce(&message, deliveryChan)
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

	bus.TaskProducer.Flush(1000)
}
