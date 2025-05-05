package usersignup

import (
	"authservice/internal/di"
	"authservice/internal/entity"
	"authservice/pkg"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type SignUpProducer struct {
	producer di.KafkaProducer
	logger   di.LoggerType
	cfg      di.ConfigType
}

func New(cfg di.ConfigType, logger di.LoggerType, producer di.KafkaProducer) *SignUpProducer {
	return &SignUpProducer{
		producer: producer,
		logger:   logger,
		cfg:      cfg,
	}
}

func (p *SignUpProducer) Close() {
	p.producer.Close()
}

func (p *SignUpProducer) SendUserSignUpInfo(user entity.UserSignUpOutDTO) {
	p.logger.Info("User SignUP Producer created successfully")

	message := pkg.CreateMessage(user, p.cfg.Kafka.UserSignupSend, -1, p.logger)

	deliveryChan := make(chan kafka.Event)

	err := p.producer.Produce(&message, deliveryChan)
	if err != nil {
		p.logger.Fatal("Failed to produce message: " + err.Error())
	}

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
