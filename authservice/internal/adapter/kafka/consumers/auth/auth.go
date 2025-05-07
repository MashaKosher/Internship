package consumers

import (
	kafkaRepo "authservice/internal/adapter/kafka"
	"authservice/internal/di"
	"authservice/internal/entity"
	"authservice/pkg"
	"fmt"
)

type AuthConsumer struct {
	consumer     di.KafkaConsumer
	logger       di.LoggerType
	cfg          di.ConfigType
	authProducer kafkaRepo.AuthProducer
	u            di.AuthService
}

func New(cfg di.ConfigType, logger di.LoggerType, consumer di.KafkaConsumer, authProducer kafkaRepo.AuthProducer, u di.AuthService) *AuthConsumer {
	return &AuthConsumer{
		consumer:     consumer,
		logger:       logger,
		cfg:          cfg,
		authProducer: authProducer,
		u:            u,
	}
}

func (p *AuthConsumer) Close() {
	p.consumer.Close()
}

func (p *AuthConsumer) ConsumerAnswerTokens() {

	p.logger.Info("Kafka Answer Auth Tokens Consumer connected successfully")
	if err := p.consumer.SubscribeTopics([]string{p.cfg.Kafka.TopicRecieve}, nil); err != nil {
		p.logger.Fatal("Failed to subscribe to topics: " + err.Error())
	}

	var authRequest entity.AuthRequest

	for {
		msg, err := p.consumer.ReadMessage(-1)
		if err == nil {

			authRequest, err := pkg.DeserializeAuthAnswer(msg.Value, authRequest, p.logger)

			if err != nil {
				p.logger.Fatal("Error while consuming: " + err.Error())
			}

			p.logger.Info("Received message: " + string(msg.Value) + " from topic:" + fmt.Sprintln(msg.TopicPartition))

			go brokerCheck(authRequest, p.u, p.logger, p.authProducer)
		} else {
			p.logger.Fatal("Error while consuming: " + err.Error())
		}
	}
}

func brokerCheck(authRequest entity.AuthRequest, u di.AuthService, logger di.LoggerType, authProducer kafkaRepo.AuthProducer) {
	var Answer entity.AuthAnswer

	user, err := u.CheckAccessToken(authRequest.AccessToken)
	if err == nil {
		Answer.Role = string(user.UserRole)
		Answer.ID = int32(user.UserID)
		Answer.Login = string(user.UserName)
		logger.Info("Token is valid, " + fmt.Sprintln(Answer))
		go authProducer.AnswerToken(Answer, authRequest.Partition)
		return
	}

	user, err = u.CheckRefreshToken(authRequest.RefreshToken)
	if err != nil {
		Answer.Err = err.Error()
		go authProducer.AnswerToken(Answer, authRequest.Partition)
		return
	}

	Answer.Role = string(user.UserRole)
	Answer.ID = int32(user.UserID)
	Answer.Login = string(user.UserName)
	Answer.NewAccessToken = user.AccessToken

	go authProducer.AnswerToken(Answer, authRequest.Partition)
}
