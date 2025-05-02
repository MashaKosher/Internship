package setup

import (
	dailyTaskRepo "coreservice/internal/adapter/db/postgres/daily_task"
	leaderboardRepo "coreservice/internal/adapter/db/postgres/leaderboard"
	seasonRepo "coreservice/internal/adapter/db/postgres/season"
	userRepo "coreservice/internal/adapter/db/postgres/user"
	seasonStatusElasticRepo "coreservice/internal/adapter/elastic/seasons"

	authCon "coreservice/internal/adapter/kafka/consumers/auth"
	dailyTaskCon "coreservice/internal/adapter/kafka/consumers/daily_task"
	matchInfoCon "coreservice/internal/adapter/kafka/consumers/match_info"
	seasonInfoCon "coreservice/internal/adapter/kafka/consumers/season_info"
	authProd "coreservice/internal/adapter/kafka/producers/auth"
	"coreservice/internal/di"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const groupID = "adminservice"
const offsetSettings = "earliest"

func mustBus(cfg di.ConfigType, logger di.LoggerType, db di.DBType, elastic di.ElasticType) di.Bus {

	authConsumer := authCon.New(cfg, logger, createConsumer(cfg, logger))
	authProducer := authProd.New(cfg, logger, createProducer(cfg, logger))
	dailyTaskConsumer := dailyTaskCon.New(cfg, logger, createConsumer(cfg, logger), dailyTaskRepo.New(db))
	matchInfoConsumer := matchInfoCon.New(
		cfg,
		logger,
		createConsumer(cfg, logger),
		userRepo.New(db),
		leaderboardRepo.New(db),
		seasonStatusElasticRepo.New(elastic.ESClient, elastic.SeasonSearchIndex, logger),
	)

	seasonInfoConsumer := seasonInfoCon.New(
		cfg,
		logger,
		createConsumer(cfg, logger),
		seasonRepo.New(db),
		seasonStatusElasticRepo.New(elastic.ESClient, elastic.SeasonSearchIndex, logger),
	)

	return di.Bus{
		AuthConsumer:       authConsumer,
		AuthProducer:       authProducer,
		DailyTaskConsumer:  dailyTaskConsumer,
		MatchInfoConsumer:  matchInfoConsumer,
		SeasonInfoConsumer: seasonInfoConsumer,
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
	bus.DailyTaskConsumer.Close()
	bus.MatchInfoConsumer.Close()
	bus.SeasonInfoConsumer.Close()
}
