package pkg

import (
	"coreservice/internal/di"
	"coreservice/internal/entity"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func serializeAuthRequest(request entity.AuthRequest, logger di.LoggerType) []byte {

	value, err := json.Marshal(request)
	if err != nil {
		logger.Fatal("Error marshaling answer: " + err.Error())
	}

	return value
}

func DeserializeAuthAnswer(value []byte, answer entity.AuthAnswer, logger di.LoggerType) (entity.AuthAnswer, error) {

	err := json.Unmarshal(value, &answer)
	logger.Info("Request recieved: " + fmt.Sprintln(answer))
	if err != nil {
		logger.Fatal("Error while consuming: " + err.Error())
	}

	return answer, err
}

func DeseriSeasonAnswer(value []byte, season entity.Season, logger di.LoggerType) (entity.Season, error) {

	err := json.Unmarshal(value, &season)
	logger.Info("Request recieved: " + fmt.Sprintln(season))
	if err != nil {
		logger.Error("Error while consuming: " + err.Error())
	}

	return season, err
}

func CreateMessage(entity entity.AuthRequest, topic string, partition int32, logger di.LoggerType) kafka.Message {

	value := serializeAuthRequest(entity, logger)
	return kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},
		Value:          value,
		Key:            []byte("a"),
	}
}
