package di

import (
	"database/sql"
	kafkaRepo "gameservice/internal/adapter/kafka"
	"gameservice/internal/config"
	"io"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
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
	Tracer     JaegerType
}

type (
	ConfigType     = config.Config
	LoggerType     = *zap.Logger
	LoggerFileType = *os.File
	DBType         = *sql.DB
	CacheType      = *redis.Client
	KafkaProducer  = *kafka.Producer
	KafkaConsumer  = *kafka.Consumer
	TracerType     = opentracing.Tracer
)

type Bus struct {
	AuthConsumer         kafkaRepo.AuthConsumer
	AuthProducer         kafkaRepo.AuthProducer
	GameSettingsConsumer kafkaRepo.GameSettingsConsumer
	MatchInfoProducer    kafkaRepo.MatchInfoProducer
}

type JaegerType struct {
	Tracer TracerType
	Closer io.Closer
}
