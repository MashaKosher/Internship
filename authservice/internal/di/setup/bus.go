package setup

import (
	"authservice/internal/di"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	authCon "authservice/internal/adapter/kafka/consumers/auth"
	authProd "authservice/internal/adapter/kafka/producers/auth"
	userSignupProd "authservice/internal/adapter/kafka/producers/user_signup"
)

func mustBus(cfg di.ConfigType, logger di.LoggerType, db di.DBType, RSAKeys di.RSAKeys) di.Bus {

	authProducer := authProd.New(cfg, logger, createProducer(cfg, logger))
	signUpProducer := userSignupProd.New(cfg, logger, createProducer(cfg, logger))
	authConsumer := authCon.New(cfg, logger, createConsumer(cfg, logger), authProducer, createAuthUseCase(db, logger, RSAKeys, signUpProducer))

	return di.Bus{
		AuthProducer:   authProducer,
		AuthConsumer:   authConsumer,
		SignUpProducer: signUpProducer,
	}
}

func createConsumer(cfg di.ConfigType, logger di.LoggerType) *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Kafka.Host + ":" + cfg.Kafka.Port,
		// "security.protocol": "PLAINTEXT",
		"group.id":          "authService",
		"auto.offset.reset": "earliest",
		// "session.timeout.ms":    10000,
		// "heartbeat.interval.ms": 3000,
		// "bootstrap.servers":        cfg.Kafka.Host + ":" + cfg.Kafka.Port,
		// "group.id":                 "authService",
		// "auto.offset.reset":        "earliest",
		// "enable.auto.commit":       false,
		// "go.events.channel.enable": true, // Используем channel-based API
		// "session.timeout.ms":       30000,
		// "heartbeat.interval.ms":    3000,
		// "max.poll.interval.ms":     300000,
		// "socket.keepalive.enable":  true,
	})
	if err != nil {
		logger.Fatal("Failed to create consumer: " + err.Error())
	}
	logger.Info("Consumer created successfully")
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
}
