package consumers

import (
	"authservice/internal/di"
	"authservice/internal/entity"
	"authservice/pkg"
	"fmt"

	producer "authservice/internal/adapter/kafka/producers"
)

func ConsumerAnswerTokens(u di.AuthService, logger di.LoggerType, cfg di.ConfigType, bus di.Bus) {
	logger.Info("Kafka Answer Auth Tokens Consumer connected successfully")

	if err := bus.Consumer.SubscribeTopics([]string{cfg.Kafka.TopicRecieve}, nil); err != nil {
		logger.Fatal("Failed to subscribe to topics: " + err.Error())
	}

	var authRequest entity.AuthRequest

	for {
		msg, err := bus.Consumer.ReadMessage(-1)
		if err == nil {

			authRequest, err := pkg.DeserializeAuthAnswer(msg.Value, authRequest, logger)

			if err != nil {
				logger.Fatal("Error while consuming: " + err.Error())
			}

			logger.Info("Received message: " + string(msg.Value) + " from topic:" + fmt.Sprintln(msg.TopicPartition))

			go BrokerCheck(authRequest, u, logger, cfg, bus)
		} else {
			logger.Fatal("Error while consuming: " + err.Error())
		}
	}
}

func BrokerCheck(authRequest entity.AuthRequest, u di.AuthService, logger di.LoggerType, cfg di.ConfigType, bus di.Bus) {
	var Answer entity.AuthAnswer

	user, err := u.CheckAccessToken(authRequest.AccessToken)
	if err != nil {
		Answer.Role = string(user.UserRole)
		Answer.ID = int32(user.UserID)
		Answer.Login = string(user.UserName)
		logger.Info("Token is valid, " + fmt.Sprintln(Answer))
		go producer.AnswerToken(Answer, authRequest.Partition, logger, cfg, bus)
		return
	}

	user, err = u.CheckRefreshToken(authRequest.RefreshToken)
	if err != nil {
		Answer.Err = err.Error()
		go producer.AnswerToken(Answer, authRequest.Partition, logger, cfg, bus)
		return
	}

	Answer.Role = string(user.UserRole)
	Answer.ID = int32(user.UserID)
	Answer.Login = string(user.UserName)
	Answer.NewAccessToken = user.AccessToken

	go producer.AnswerToken(Answer, authRequest.Partition, logger, cfg, bus)
}
