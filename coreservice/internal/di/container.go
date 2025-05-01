package di

import (
	"coreservice/internal/config"
	db "coreservice/internal/repository/sqlc/generated"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Container struct {
	Config     ConfigType
	Logger     LoggerType
	LoggerFile LoggerFileType
	Services   Services
	// Bus        Bus
	DB        DBType
	Validator ValidatorType
	Elastic   ElasticType
}

type (
	ConfigType     = config.Config
	LoggerType     = *zap.Logger
	LoggerFileType = *os.File
	DBType         = *db.Queries
	ValidatorType  = *validator.Validate
	ElasticIndex   = string
	ESClient       = *elasticsearch.Client
)

type ElasticType struct {
	ESClient          ESClient
	UserSearchIndex   ElasticIndex
	SeasonSearchIndex ElasticIndex
}
