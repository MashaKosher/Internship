package consumers

import (
	"authservice/internal/config"
	"authservice/internal/entity"
	"authservice/internal/usecase"
	"authservice/pkg"
	"authservice/pkg/logger"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	producer "authservice/internal/adapter/kafka/producers"
)

func ConsumerAnswerTokens(u usecase.Auth) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.AppConfig.Kafka.Host + ":" + config.AppConfig.Kafka.Port,
		"group.id":          "authService",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		logger.Logger.Fatal("Failed to create consumer: " + err.Error())
	}
	defer c.Close()

	logger.Logger.Info("Kafka Answer Auth Tokens Consumer connected successfully")

	if err := c.SubscribeTopics([]string{config.AppConfig.Kafka.TopicRecieve}, nil); err != nil {
		logger.Logger.Fatal("Failed to subscribe to topics: " + err.Error())
	}

	var authRequest entity.AuthRequest

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {

			// Тут поменять и вообще все кафка утилс
			authRequest, err := pkg.DeserializeAuthAnswer(msg.Value, authRequest)

			if err != nil {
				logger.Logger.Fatal("Error while consuming: " + err.Error())
			}

			logger.Logger.Info("Received message: " + string(msg.Value) + " from topic:" + fmt.Sprintln(msg.TopicPartition))

			// if we recieve Auth request
			go BrokerCheck(authRequest, u)
		} else {
			logger.Logger.Fatal("Error while consuming: " + err.Error())
		}
	}
}

func BrokerCheck(authRequest entity.AuthRequest, u usecase.Auth) {
	var Answer entity.AuthAnswer

	user, err := u.CheckAccessToken(authRequest.AccessToken)
	if err != nil {
		Answer.Role = string(user.UserRole)
		Answer.ID = int32(user.UserID)
		Answer.Login = string(user.UserName)
		logger.Logger.Info("Token is valid, " + fmt.Sprintln(Answer))
		go producer.AnswerToken(Answer, authRequest.Partition)
		return
	}

	user, err = u.CheckRefreshToken(authRequest.RefreshToken)
	if err != nil {
		Answer.Err = err.Error()
		go producer.AnswerToken(Answer, authRequest.Partition)
		return
	}

	Answer.Role = string(user.UserRole)
	Answer.ID = int32(user.UserID)
	Answer.Login = string(user.UserName)
	Answer.NewAccessToken = user.AccessToken

	go producer.AnswerToken(Answer, authRequest.Partition)
}
