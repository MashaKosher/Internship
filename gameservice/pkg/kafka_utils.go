package pkg

import (
	"encoding/json"
	"fmt"
	"gameservice/internal/entity"
	"gameservice/pkg/logger"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func CreateMessage(entity any, topic string, partition int32) kafka.Message {
	value := serializeRequest(entity)
	return kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},
		Value:          value,
		Key:            []byte("a"),
	}
}

func serializeRequest(request any) []byte {
	value, err := json.Marshal(request)
	if err != nil {
		logger.L.Fatal("Error marshaling answer: " + err.Error())
	}
	return value
}

func DeserializeAuthAnswer(value []byte, answer entity.AuthAnswer) (entity.AuthAnswer, error) {

	err := json.Unmarshal(value, &answer)
	logger.L.Info("Request recieved: " + fmt.Sprintln(answer))
	if err != nil {
		logger.L.Fatal("Error while consuming: " + err.Error())
	}

	return answer, err
}
