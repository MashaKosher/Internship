package auth

import (
	"adminservice/internal/di"
	"adminservice/internal/entity"
	utils "adminservice/pkg/kafka_utils"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type AuthConsumer struct {
	consumer di.KafkaConsumer
	logger   di.LoggerType
	cfg      di.ConfigType
}

func New(cfg di.ConfigType, logger di.LoggerType, consumer di.KafkaConsumer) *AuthConsumer {
	return &AuthConsumer{
		consumer: consumer,
		logger:   logger,
		cfg:      cfg,
	}
}

func (c *AuthConsumer) Close() {
	c.consumer.Close()
}

func (c *AuthConsumer) AnswerTokens() (entity.AuthAnswer, error) {

	var err error
	var answer entity.AuthAnswer

	c.logger.Info("Kafka Auth consumer connected successfully")

	err = c.consumer.Assign([]kafka.TopicPartition{{Topic: &c.cfg.Kafka.AuthTopicSend, Partition: c.cfg.Kafka.Partition, Offset: kafka.OffsetTail(1)}})
	if err != nil {
		c.logger.Fatal("Failed to assign partition:" + err.Error())
	}

	for {
		msg, err := c.consumer.ReadMessage(-1)
		if err == nil {
			c.logger.Info("Received message: " + string(msg.Value) + " from topic:" + msg.TopicPartition.String() + " with offset " + msg.TopicPartition.Offset.String())

			answer, err := utils.DeserializeAuthAnswer(msg.Value, answer)
			if err != nil {
				c.logger.Error("Error while consuming: " + err.Error())
				return answer, err
			}
			return answer, err
		} else {
			c.logger.Error("Error while consuming: " + err.Error())
			return answer, err
		}
	}
}
