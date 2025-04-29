package consumers

import (
	"gameservice/internal/di"
	"gameservice/internal/entity"
	"gameservice/pkg"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func RecieveTokenInfo(cfg di.ConfigType, bus di.Bus) (entity.AuthAnswer, error) {

	var err error
	var answer entity.AuthAnswer

	// consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
	// 	"bootstrap.servers": cfg.Kafka.Host + ":" + cfg.Kafka.Port, // Используйте localhost
	// 	"group.id":          "adminServiceBecomingAnswer",
	// 	"auto.offset.reset": "latest",
	// })
	// if err != nil {
	// 	bus.Logger.Fatal("Failed to create consumer:" + err.Error())
	// }
	// defer consumer.Close()
	bus.Logger.Info("Auth Consumer connected successfully")

	err = bus.AuthConsumer.Assign([]kafka.TopicPartition{{Topic: &cfg.Kafka.AuthTopicSend, Partition: cfg.Kafka.Partition, Offset: kafka.OffsetTail(1)}})
	if err != nil {
		bus.Logger.Fatal("Failed to assign partition:" + err.Error())
	}

	for {
		msg, err := bus.AuthConsumer.ReadMessage(-1)
		if err == nil {
			bus.Logger.Info("Received message: " + string(msg.Value) + " from topic:" + msg.TopicPartition.String() + " with offset " + msg.TopicPartition.Offset.String())

			answer, err := pkg.DeserializeAuthAnswer(msg.Value, answer, bus.Logger)
			if err != nil {
				bus.Logger.Error("Error while consuming:" + err.Error())
				return answer, err
			}
			return answer, err
		} else {
			bus.Logger.Error("Error while consuming: " + err.Error())
			return answer, err
		}
	}
}
