package consumers

import (
	"adminservice/internal/entity"
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func AnswerTokens() (entity.AuthAnswer, error) {

	var err error
	var answer entity.AuthAnswer

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092", // Используйте localhost
		// "bootstrap.servers": "kafka:9092", // Используйте localhost
		"group.id":          "adminServiceBecomingAnswer",
		"auto.offset.reset": "latest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer consumer.Close()
	log.Println("Kafka connected successfully")

	topic := "jwtCheckAnswer"

	// if err := consumer.SubscribeTopics([]string{"jwtCheckAnswer"}, nil); err != nil {
	// 	log.Fatal("Failed to subscribe to topics: " + err.Error())
	// 	panic(err)
	// }

	err = consumer.Assign([]kafka.TopicPartition{{Topic: &topic, Partition: 0, Offset: kafka.OffsetTail(1)}})

	if err != nil {
		log.Fatal("Failed to assign partition:", err)
	}

	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			log.Println("Received message: " + string(msg.Value) + " from topic:" + msg.TopicPartition.String() + " with offset " + msg.TopicPartition.Offset.String())
			err = json.Unmarshal(msg.Value, &answer)
			log.Println("Request recieved: " + fmt.Sprintln(answer))
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
