package consumers

import (
	"coreservice/internal/di"
	"coreservice/internal/entity"
	"coreservice/pkg"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func RecieveTokenInfo(cfg di.ConfigType, bus di.Bus) (entity.AuthAnswer, error) {
	var err error
	var answer entity.AuthAnswer

	// consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
	// 	"bootstrap.servers": config.AppConfig.Kafka.Host + ":" + config.AppConfig.Kafka.Port, // Используйте localhost
	// 	"group.id":          "authService",
	// 	"auto.offset.reset": "earliest",
	// })
	// if err != nil {
	// 	logger.Logger.Error("Failed to create consumer: " + err.Error())
	// }
	// defer consumer.Close()

	bus.Logger.Info("Kafka connected successdully")
	err = bus.AuthConsumer.Assign([]kafka.TopicPartition{{Topic: &cfg.Kafka.AuthTopicRecieve, Partition: cfg.Kafka.Partition, Offset: kafka.OffsetTail(1)}})
	if err != nil {
		bus.Logger.Fatal("Failed to assign partition:" + err.Error())
	}

	for {
		msg, err := bus.AuthConsumer.ReadMessage(-1)
		if err == nil {
			bus.Logger.Info("Received message: " + string(msg.Value) + " from topic:" + msg.TopicPartition.String() + " with offset " + msg.TopicPartition.Offset.String())

			answer, err := pkg.DeserializeAuthAnswer(msg.Value, answer, bus.Logger)
			if err != nil {
				bus.Logger.Error("Error while consuming: " + err.Error())
				return answer, err
			}
			return answer, err
		} else {
			bus.Logger.Error("Error while consuming: " + err.Error())
			return answer, err
		}
	}
}
