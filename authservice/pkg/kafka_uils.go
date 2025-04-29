package pkg

import (
	"authservice/internal/di"
	"authservice/internal/entity"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func CreateMessage(entity entity.AuthAnswer, topic string, partition int32, logger di.LoggerType) kafka.Message {
	value := serializeAuthAnswer(entity, logger)
	return kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},
		Value:          value,
		Key:            []byte("a"),
	}
}

func serializeAuthAnswer(answer entity.AuthAnswer, logger di.LoggerType) []byte {
	value, err := json.Marshal(answer)
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
