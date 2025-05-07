package di

import (
	"adminservice/internal/config"
	"io"
	"os"

	kafkaRepo "adminservice/internal/adapter/kafka"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Container struct {
	Config     ConfigType
	Logger     LoggerType
	LoggerFile LoggerFileType
	Services   Services
	Bus        Bus
	DB         DBType
	Validator  ValidatorType
	Tracer     JaegerType
}

type (
	ConfigType     = config.Config
	LoggerType     = *zap.Logger
	LoggerFileType = *os.File
	DBType         = *gorm.DB
	ValidatorType  = *validator.Validate
	KafkaProducer  = *kafka.Producer
	KafkaConsumer  = *kafka.Consumer
	TracerType     = opentracing.Tracer
)

type Bus struct {
	AuthConsumer         kafkaRepo.AuthConsumer
	AuthProducer         kafkaRepo.AuthProducer
	GameSettingsProducer kafkaRepo.GameSettingsProducer
	SeasonProducer       kafkaRepo.SeasonProducer
	DailyTaskProducer    kafkaRepo.DailyTaskProducer
}

type JaegerType struct {
	Tracer TracerType
	Closer io.Closer
}
