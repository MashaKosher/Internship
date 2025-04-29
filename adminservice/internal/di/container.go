package di

import (
	"adminservice/internal/config"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-playground/validator/v10"
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
}

// ///////////////////////////////////
type (
	ConfigType     = config.Config
	LoggerType     = *zap.Logger
	LoggerFileType = *os.File
	DBType         = *gorm.DB
	ValidatorType  = *validator.Validate
)

type Bus struct {
	Consumer       *kafka.Consumer
	AuthProducer   *kafka.Producer
	GameProducer   *kafka.Producer
	SeasonProducer *kafka.Producer
	TaskProducer   *kafka.Producer
	Logger         LoggerType
}
