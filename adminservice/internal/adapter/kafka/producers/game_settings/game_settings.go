package game

import (
	"adminservice/internal/di"
	"adminservice/internal/entity"
	utils "adminservice/pkg/kafka_utils"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type GameSettingsProducer struct {
	producer di.KafkaProducer
	logger   di.LoggerType
	cfg      di.ConfigType
}

func New(cfg di.ConfigType, logger di.LoggerType, producer di.KafkaProducer) *GameSettingsProducer {
	return &GameSettingsProducer{
		producer: producer,
		logger:   logger,
		cfg:      cfg,
	}
}

func (p *GameSettingsProducer) Close() {
	p.producer.Close()
}

func (p *GameSettingsProducer) SendGameSettings(season entity.SettingsJson) {
	p.logger.Info("Season Producer created successfully")

	message := utils.CreateMessage(season, p.cfg.Kafka.GameSettingsTopicSend, -1, p.logger)

	deliveryChan := make(chan kafka.Event)

	err := p.producer.Produce(&message, deliveryChan)
	if err != nil {
		p.logger.Error("Failed to produce message: " + err.Error())
		return
	}

	p.logger.Info("Message sent, waiting for delivery confirmation...")

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
			p.logger.Info("Ignored event: " + e.String())
		}
	}()

	p.producer.Flush(1000)
}
