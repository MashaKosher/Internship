package pkg

import (
	"authservice/internal/entity"
	"authservice/internal/logger"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func CreateMessage(entity entity.AuthAnswer, topic string, partition int32) kafka.Message {

	value := serializeAuthAnswer(entity)
	return kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},
		Value:          value,
		Key:            []byte("a"),
	}
}

func serializeAuthAnswer(answer entity.AuthAnswer) []byte {
	value, err := json.Marshal(answer)
	if err != nil {
		logger.Logger.Fatal("Error marshaling answer: " + err.Error())
	}

	return value
}

func DeserializeAuthAnswer(value []byte, request entity.AuthRequest) (entity.AuthRequest, error) {
	err := json.Unmarshal(value, &request)
	logger.Logger.Info("Request recieved: " + fmt.Sprintln(request))
	if err != nil {
		logger.Logger.Error("Error while consuming: " + err.Error())
	}

	return request, err
}
