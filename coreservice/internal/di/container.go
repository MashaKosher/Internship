package di

import (
	kafkaRepo "coreservice/internal/adapter/kafka"
	"coreservice/internal/config"
	db "coreservice/internal/repository/sqlc/generated"
	"io"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"

	elasticRepo "coreservice/internal/adapter/elastic"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type Container struct {
	Config     ConfigType
	Logger     LoggerType
	LoggerFile LoggerFileType
	Services   Services
	Bus        Bus
	DB         DBType
	Validator  ValidatorType
	Elastic    ElasticType
	Cache      CacheType
	Tracer     JaegerType
}

type (
	ConfigType     = config.Config
	LoggerType     = *zap.Logger
	LoggerFileType = *os.File
	DBType         = *db.Queries
	ValidatorType  = *validator.Validate
	ElasticIndex   = string
	ESClient       = *elasticsearch.Client
	DelayProducer  = *asynq.Client
	CacheType      = *redis.Client
	KafkaProducer  = *kafka.Producer
	KafkaConsumer  = *kafka.Consumer
	TracerType     = opentracing.Tracer
)

type Bus struct {
	AuthConsumer       kafkaRepo.AuthConsumer
	AuthProducer       kafkaRepo.AuthProducer
	DailyTaskConsumer  kafkaRepo.DailyTaskConsumer
	MatchInfoConsumer  kafkaRepo.MatchInfoConsumer
	SeasonInfoConsumer kafkaRepo.SeasonInfoConsumer
	UserSignUpConsumer kafkaRepo.UserSignUpConsumer
}

type ElasticType struct {
	SeasonStatus elasticRepo.SeasonStatusRepo
	UserName     elasticRepo.UserNameRepo
}

type JaegerType struct {
	Tracer TracerType
	Closer io.Closer
}
