package consumers

import (
	"adminservice/internal/di"
	"adminservice/internal/entity"
	"adminservice/pkg"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func AnswerTokens(cfg di.ConfigType, bus di.Bus) (entity.AuthAnswer, error) {

	var err error
	var answer entity.AuthAnswer

	// consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
	// 	"bootstrap.servers": config.AppConfig.Kafka.Host + ":" + config.AppConfig.Kafka.Port, // Используйте localhost
	// 	"group.id":          "adminServiceBecomingAnswer",
	// 	"auto.offset.reset": "latest",
	// })
	// if err != nil {
	// 	log.Fatalf("Failed to create consumer: %s", err)
	// }
	// defer consumer.Close()

	log.Println("Kafka connected successfully")

	err = bus.Consumer.Assign([]kafka.TopicPartition{{Topic: &cfg.Kafka.AuthTopicSend, Partition: cfg.Kafka.Partition, Offset: kafka.OffsetTail(1)}})
	if err != nil {
		log.Fatal("Failed to assign partition:", err)
	}

	for {
		msg, err := bus.Consumer.ReadMessage(-1)
		if err == nil {
			log.Println("Received message: " + string(msg.Value) + " from topic:" + msg.TopicPartition.String() + " with offset " + msg.TopicPartition.Offset.String())

			answer, err := pkg.DeserializeAuthAnswer(msg.Value, answer)
			if err != nil {
				log.Printf("Error while consuming: %v", err)
				return answer, err
			}
			return answer, err
		} else {
			log.Printf("Error while consuming: %v", err)
			return answer, err
		}
	}
}
