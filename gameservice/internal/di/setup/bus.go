package setup

import (
	"gameservice/internal/di"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func mustBus(cfg di.ConfigType, logger di.LoggerType) di.Bus {
	return di.Bus{
		AuthConsumer:         createConsumer(cfg, logger),
		GameSettingsConsumer: createConsumer(cfg, logger),
		AuthProducer:         createProducer(cfg, logger),
		Logger:               logger,
	}
}

func createConsumer(cfg di.ConfigType, logger di.LoggerType) *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.Host + ":" + cfg.Kafka.Port,
		"group.id":          "gameservice",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		logger.Fatal("Failed to create consumer: " + err.Error())
	}
	return c
}

func createProducer(cfg di.ConfigType, logger di.LoggerType) *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.Kafka.Host + ":" + cfg.Kafka.Port})
	if err != nil {
		logger.Error("Failed to create producer:" + err.Error())
	}

	return p
}

func deferBus(bus di.Bus) {
	bus.AuthConsumer.Close()
	bus.GameSettingsConsumer.Close()
	bus.AuthProducer.Close()
}
