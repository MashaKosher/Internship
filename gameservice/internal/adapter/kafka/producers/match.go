package producers

import (
	"fmt"
	"gameservice/internal/di"
	"gameservice/internal/entity"
	"gameservice/pkg"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func SendMatchInfo(match entity.GameResult, cfg di.ConfigType, bus di.Bus) {
	// p, err := kafka.NewProducer(&kafka.ConfigMap{
	// 	"bootstrap.servers": cfg.Kafka.Host + ":" + cfg.Kafka.Port,
	// })
	// if err != nil {
	// 	bus.Logger.Fatal("Failed to create producer: " + err.Error())
	// }
	// defer p.Close()

	bus.Logger.Info("Match Producer created successfully")

	message := pkg.CreateMessage(match, cfg.Kafka.MatchTopicSend, -1, bus.Logger)

	// Канал для получения событий доставки
	deliveryChan := make(chan kafka.Event)

	err := bus.MatchProducer.Produce(&message, deliveryChan)
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
			bus.Logger.Info("Ignored event: " + e.String())
		}
	}()

	bus.MatchProducer.Flush(1000)
}
