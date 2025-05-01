package di

import (
	"coreservice/internal/config"
	db "coreservice/internal/repository/sqlc/generated"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-playground/validator/v10"
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
	// DelayTask  DelayTaskType
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
	// DelayConsumer  = *asynq.Server
)

type ElasticType struct {
	ESClient          ESClient
	UserSearchIndex   ElasticIndex
	SeasonSearchIndex ElasticIndex
}

type Bus struct {
	AuthConsumer      *kafka.Consumer
	DailyTaskConsumer *kafka.Consumer
	SeasonConsumer    *kafka.Consumer
	AuthProducer      *kafka.Producer
	Logger            LoggerType
}

// type DelayTaskType struct {
// 	SeasonConsumer *asynq.Server
// 	SeasonProducer *asynq.Client
// 	Logger         LoggerType
// }
