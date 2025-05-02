package setup

import (
	authCons "adminservice/internal/adapter/kafka/consumers/auth"
	authProds "adminservice/internal/adapter/kafka/producers/auth"
	dailutaskProds "adminservice/internal/adapter/kafka/producers/daily_task"
	gamesettingsProds "adminservice/internal/adapter/kafka/producers/game_settings"
	seasonProds "adminservice/internal/adapter/kafka/producers/season"
	"adminservice/internal/di"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func mustBus(cfg di.ConfigType, logger di.LoggerType) di.Bus {

	authConsumer := authCons.New(cfg, logger, createConsumer(cfg, logger))
	authProducer := authProds.New(cfg, logger, createProducer(cfg, logger))
	gameSettingsProducer := gamesettingsProds.New(cfg, logger, createProducer(cfg, logger))
	seasonProducer := seasonProds.New(cfg, logger, createProducer(cfg, logger))
	dailyTaskProducer := dailutaskProds.New(cfg, logger, createProducer(cfg, logger))

	return di.Bus{
		AuthConsumer:         authConsumer,
		AuthProducer:         authProducer,
		GameSettingsProducer: gameSettingsProducer,
		SeasonProducer:       seasonProducer,
		DailyTaskProducer:    dailyTaskProducer,
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
	bus.AuthConsumer.Close()
	bus.AuthProducer.Close()
	bus.GameSettingsProducer.Close()
	bus.SeasonProducer.Close()
	bus.DailyTaskProducer.Close()
}
