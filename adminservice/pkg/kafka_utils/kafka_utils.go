package kafkautils

import (
	"adminservice/internal/di"
	"adminservice/internal/entity"
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func CreateMessage(entity any, topic string, partition int32, logger di.LoggerType) kafka.Message {
	value := serializeEntity(entity, logger)
	return kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},
		Value:          value,
		Key:            nil,
	}
}

func serializeEntity(request any, logger di.LoggerType) []byte {

	value, err := json.Marshal(request)
	if err != nil {
		logger.Fatal("Error marshaling answer: " + err.Error())
	}

	return value
}

func DeserializeAuthAnswer(value []byte, answer entity.AuthAnswer) (entity.AuthAnswer, error) {

	err := json.Unmarshal(value, &answer)
	log.Println("Request recieved: " + fmt.Sprintln(answer))
	if err != nil {
		log.Fatal("Error while consuming: " + err.Error())
	}

	return answer, err
}
