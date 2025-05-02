package pkg

import (
	"encoding/json"
	"fmt"
	"gameservice/internal/di"
	"gameservice/internal/entity"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func CreateMessage(entity any, topic string, partition int32, logger di.LoggerType) kafka.Message {
	value := serializeRequest(entity, logger)
	return kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: partition},
		Value:          value,
		Key:            nil,
	}
}

func serializeRequest(request any, logger di.LoggerType) []byte {
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

func DeserializeGameSettings(value []byte, logger di.LoggerType) (entity.GameSettings, error) {

	var gameSettings entity.GameSettings
	err := json.Unmarshal(value, &gameSettings)
	logger.Info("Game Settings recieved: " + fmt.Sprintln(gameSettings))
	if err != nil {
		logger.Fatal("Error while consuming: " + err.Error())
	}

	return gameSettings, err
}
