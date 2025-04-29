package di

import (
	"authservice/internal/config"
	"crypto/rsa"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TokenType string

const ACCESS_TOKEN TokenType = "access"
const REFRESH_TOKEN TokenType = "refresh"

type Container struct {
	Config     ConfigType
	Logger     LoggerType
	LoggerFile LoggerFileType
	Services   Services
	Bus        Bus
	DB         DBType
	RSAKeys    RSAKeys
	Validator  ValidatorType
}

// //////////////////////////////////
type (
	ConfigType     = config.Config
	LoggerType     = *zap.Logger
	DBType         = *gorm.DB
	LoggerFileType = *os.File
	ValidatorType  = *validator.Validate
)

type Bus struct {
	Consumer *kafka.Consumer
	Producer *kafka.Producer
	Logger   LoggerType
}

type RSAKeys struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}
