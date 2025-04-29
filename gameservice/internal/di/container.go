package di

import (
	"database/sql"
	"gameservice/internal/config"
	"os"

	"github.com/redis/go-redis/v9"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type Container struct {
	Config     ConfigType
	Logger     LoggerType
	LoggerFile LoggerFileType
	Services   Services
	Bus        Bus
	DB         DBType
	Cache      CacheType
}

// ///////////////////////////////////
type (
	ConfigType     = config.Config
	LoggerType     = *zap.Logger
	LoggerFileType = *os.File
	DBType         = *sql.DB
	CacheType      = *redis.Client
)

type Bus struct {
	AuthConsumer         *kafka.Consumer
	GameSettingsConsumer *kafka.Consumer
	AuthProducer         *kafka.Producer
	Logger               LoggerType
}
