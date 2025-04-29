package setup

import (
	"adminservice/internal/di"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func mustBus(cfg di.ConfigType, logger di.LoggerType) di.Bus {
	consumer := createConsumer(cfg, logger)

	return di.Bus{
		Consumer:       consumer,
		AuthProducer:   createProducer(cfg, logger),
		GameProducer:   createProducer(cfg, logger),
		SeasonProducer: createProducer(cfg, logger),
		TaskProducer:   createProducer(cfg, logger),
		Logger:         logger,
	}
}

func createConsumer(cfg di.ConfigType, logger di.LoggerType) *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.Host + ":" + cfg.Kafka.Port,
		"group.id":          "adminservice",
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
	bus.Consumer.Close()
	bus.AuthProducer.Close()
	bus.GameProducer.Close()
	bus.SeasonProducer.Close()
	bus.TaskProducer.Close()
}
