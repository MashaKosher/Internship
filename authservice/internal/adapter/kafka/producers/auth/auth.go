package producers

import (
	"authservice/internal/di"
	"authservice/internal/entity"
	"authservice/pkg"
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

func (p *AuthProducer) AnswerToken(answer entity.AuthAnswer, partition int32) {
	p.logger.Info("Auth Producer created successfully")

	message := pkg.CreateMessage(answer, p.cfg.Kafka.TopicSend, partition, p.logger)

	deliveryChan := make(chan kafka.Event)

	err := p.producer.Produce(&message, deliveryChan)
	if err != nil {
		p.logger.Fatal("Failed to produce message: " + err.Error())
	}

	p.logger.Info("Produced message: " + fmt.Sprint(answer))

	go func() {
		event := <-deliveryChan
		switch e := event.(type) {
		case *kafka.Message:
			if e.TopicPartition.Error != nil {
				p.logger.Error(fmt.Sprintf("Delivery failed: %v", e.TopicPartition.Error))
			} else {
				p.logger.Error(fmt.Sprintf("Delivered message to %v [%d] at offset %v",
					*e.TopicPartition.Topic, e.TopicPartition.Partition, e.TopicPartition.Offset))
			}
		default:
			p.logger.Error(fmt.Sprintf("Ignored event: %s", e))
		}
	}()

	p.producer.Flush(1000)
}
