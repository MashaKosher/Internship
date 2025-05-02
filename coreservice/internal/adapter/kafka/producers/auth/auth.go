package auth

import (
	"coreservice/internal/di"
	"coreservice/internal/entity"
	utils "coreservice/pkg/kafka_utils"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type AuthProducer struct {
	producer di.KafkaProducer
	logger   di.LoggerType
	cfg      di.ConfigType
}

func New(cfg di.ConfigType, logger di.LoggerType, producer di.KafkaProducer) *AuthProducer {
	return &AuthProducer{
		producer: producer,
		logger:   logger,
		cfg:      cfg,
	}
}

func (p *AuthProducer) Close() {
	p.producer.Close()
}

func (p *AuthProducer) CheckAuthTokenRequest(accessToken, refreshToken string) {
	p.logger.Info("Producer created successfully")

	var request entity.AuthRequest
	request.Partition = p.cfg.Kafka.Partition
	request.AccessToken = accessToken
	request.RefreshToken = refreshToken

	message := utils.CreateMessage(request, p.cfg.Kafka.AuthTopicRecieve, p.cfg.Kafka.Partition, p.logger)

	// Канал для получения событий доставки
	deliveryChan := make(chan kafka.Event)

	err := p.producer.Produce(&message, deliveryChan)
	if err != nil {
		p.logger.Error("Failed to produce message: " + err.Error())
		return
	}

	p.logger.Info("Message sent, waiting for delivery confirmation...")

	// Ожидаем подтверждения доставки
	go func() {
		event := <-deliveryChan
		switch e := event.(type) {
		case *kafka.Message:
			if e.TopicPartition.Error != nil {
				p.logger.Error("Delivery failed: " + e.TopicPartition.Error.Error())
			} else {
				p.logger.Info(fmt.Sprintf("Delivered message to %v [%d] at offset %v",
					*e.TopicPartition.Topic, e.TopicPartition.Partition, e.TopicPartition.Offset))
			}
		default:
			p.logger.Error("Ignored event: " + e.String())
		}
	}()

	p.producer.Flush(1000)
}
