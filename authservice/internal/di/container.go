package di

import (
	kafkaRepo "authservice/internal/adapter/kafka"
	"authservice/internal/adapter/memcached"
	"authservice/internal/config"
	"crypto/rsa"
	"io"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-playground/validator/v10"
	"github.com/opentracing/opentracing-go"

	"github.com/bradfitz/gomemcache/memcache"
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
	Tracer     JaegerType
	Cache      Cache
}

// ///////////////////////////////////
type (
	ConfigType     = config.Config
	LoggerType     = *zap.Logger
	LoggerFileType = *os.File
	DBType         = *gorm.DB
	ValidatorType  = *validator.Validate
	KafkaProducer  = *kafka.Producer
	KafkaConsumer  = *kafka.Consumer
	TracerType     = opentracing.Tracer
	CacheType      = *memcache.Client
)

type Bus struct {
	AuthProducer   kafkaRepo.AuthProducer
	AuthConsumer   kafkaRepo.AuthConsumer
	SignUpProducer kafkaRepo.SignUpProducer
}

type RSAKeys struct {
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
}

type JaegerType struct {
	Tracer TracerType
	Closer io.Closer
}

type Cache struct {
	Token memcached.TokenCacheRepo
}
