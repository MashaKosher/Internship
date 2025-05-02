package setup

import (
	authCon "adminservice/internal/adapter/kafka/consumers/auth"
	authProd "adminservice/internal/adapter/kafka/producers/auth"
	dailutaskProd "adminservice/internal/adapter/kafka/producers/daily_task"
	gamesettingsProd "adminservice/internal/adapter/kafka/producers/game_settings"
	seasonProd "adminservice/internal/adapter/kafka/producers/season"
	"adminservice/internal/di"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const groupID = "adminservice"
const offsetSettings = "earliest"

func mustBus(cfg di.ConfigType, logger di.LoggerType) di.Bus {

	authConsumer := authCon.New(cfg, logger, createConsumer(cfg, logger))
	authProducer := authProd.New(cfg, logger, createProducer(cfg, logger))
	gameSettingsProducer := gamesettingsProd.New(cfg, logger, createProducer(cfg, logger))
	seasonProducer := seasonProd.New(cfg, logger, createProducer(cfg, logger))
	dailyTaskProducer := dailutaskProd.New(cfg, logger, createProducer(cfg, logger))

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
		"group.id":          groupID,
		"auto.offset.reset": offsetSettings,
	})
	if err != nil {
		logger.Fatal("Failed to create consumer: " + err.Error())
	}
	return c
}

func createProducer(cfg di.ConfigType, logger di.LoggerType) *kafka.Producer {
	p, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": cfg.Kafka.Host + ":" + cfg.Kafka.Port,
		})
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
