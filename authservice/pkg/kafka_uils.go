package pkg

import (
	"authservice/internal/di"
	"authservice/internal/entity"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func CreateMessage[T entity.AuthAnswer | entity.UserSignUpOutDTO](entity T, topic string, partition int32, logger di.LoggerType) kafka.Message {
	value := serializeEntity(entity, logger)
	return kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},
		Value:          value,
		Key:            nil,
	}
}

func serializeEntity[T entity.AuthAnswer | entity.UserSignUpOutDTO](entity T, logger di.LoggerType) []byte {
	value, err := json.Marshal(entity)
	if err != nil {
		logger.Fatal("Error marshaling answer: " + err.Error())
	}

	return value
}

func DeserializeAuthAnswer(value []byte, request entity.AuthRequest, logger di.LoggerType) (entity.AuthRequest, error) {
	err := json.Unmarshal(value, &request)
	logger.Info("Request recieved: " + fmt.Sprintln(request))
	if err != nil {
		logger.Error("Error while consuming: " + err.Error())
	}

	return request, err
}
