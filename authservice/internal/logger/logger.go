package logger

import (
	cnf "authservice/internal/config"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func CreateLogger() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout", cnf.Cfg.LogFileName}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	Logger = logger
	Logger.Info("Program start")

}
