package pkg

import (
	"coreservice/internal/entity"
	"coreservice/internal/logger"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func serializeAuthRequest(request entity.AuthRequest) []byte {

	value, err := json.Marshal(request)
	if err != nil {
		logger.Logger.Fatal("Error marshaling answer: " + err.Error())
		panic(err)
	}

	return value
}

func DeserializeAuthAnswer(value []byte, answer entity.AuthAnswer) (entity.AuthAnswer, error) {

	err := json.Unmarshal(value, &answer)
	logger.Logger.Info("Request recieved: " + fmt.Sprintln(answer))
	if err != nil {
		logger.Logger.Error("Error while consuming: " + err.Error())
	}

	return answer, err
}

func DeseriSeasonAnswer(value []byte, season entity.Season) (entity.Season, error) {

	err := json.Unmarshal(value, &season)
	logger.Logger.Info("Request recieved: " + fmt.Sprintln(season))
	if err != nil {
		logger.Logger.Error("Error while consuming: " + err.Error())
	}

	return season, err
}

func CreateMessage(entity entity.AuthRequest, topic string, partition int32) kafka.Message {

	value := serializeAuthRequest(entity)
	return kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},
		Value:          value,
		Key:            []byte("a"),
	}
}
