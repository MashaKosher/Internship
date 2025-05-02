package kafka

import "authservice/internal/entity"

type (
	AuthProducer interface {
		Close()
		AnswerToken(answer entity.AuthAnswer, partition int32)
	}

	AuthConsumer interface {
		Close()
		ConsumerAnswerTokens()
	}
)
