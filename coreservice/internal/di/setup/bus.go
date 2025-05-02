package setup

import (
	"coreservice/internal/di"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func mustBus(cfg di.ConfigType, logger di.LoggerType) di.Bus {
	return di.Bus{
		AuthConsumer:      createBusConsumer(cfg, logger),
		DailyTaskConsumer: createBusConsumer(cfg, logger),
		SeasonConsumer:    createBusConsumer(cfg, logger),
		GameConsumer:      createBusConsumer(cfg, logger),
		AuthProducer:      createBusProducer(cfg, logger),
		Logger:            logger,
	}
}

func createBusConsumer(cfg di.ConfigType, logger di.LoggerType) *kafka.Consumer {
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

func createBusProducer(cfg di.ConfigType, logger di.LoggerType) *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": cfg.Kafka.Host + ":" + cfg.Kafka.Port})
	if err != nil {
		logger.Error("Failed to create producer:" + err.Error())
	}

	return p
}

func deferBus(bus di.Bus) {
	bus.AuthConsumer.Close()
	bus.DailyTaskConsumer.Close()
	bus.SeasonConsumer.Close()
	bus.GameConsumer.Close()
	bus.AuthProducer.Close()
}
