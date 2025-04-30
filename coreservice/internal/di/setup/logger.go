package setup

import (
	"coreservice/internal/di"

	"go.uber.org/zap"
)

func mustLogger(cfg di.ConfigType) di.LoggerType {
	conf := zap.NewProductionConfig()
	conf.OutputPaths = []string{"stdout", cfg.Logger.FileName}

	logger, err := conf.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

func deferLogger(logger di.LoggerType) {
	logger.Info("Program end")
	logger.Sync()
}
