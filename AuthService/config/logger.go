package config

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func CreateLogger() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout", Envs.LogFileName}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	Logger = logger
	Logger.Info("Program start")

}
