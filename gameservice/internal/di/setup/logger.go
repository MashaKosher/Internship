package setup

import (
	"gameservice/internal/di"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func mustLogger(cfg di.ConfigType) di.LoggerType {
	conf := zap.NewProductionConfig()
	conf.OutputPaths = []string{"stdout", cfg.Logger.FileName}
	conf.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

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
