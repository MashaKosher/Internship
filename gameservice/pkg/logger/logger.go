package logger

import (
	"gameservice/internal/config"
	"os"

	"go.uber.org/zap"
)

var L *zap.Logger

func CreateLogger() *os.File {
	conf := zap.NewProductionConfig()
	conf.OutputPaths = []string{"stdout", config.AppConfig.Logger.FileName}

	logger, err := conf.Build()
	if err != nil {
		panic(err)
	}

	L = logger
	L.Info("Program started")

	return createLogFile()

}

func createLogFile() *os.File {
	logFile, err := os.OpenFile(config.AppConfig.Logger.FileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("Log error: " + err.Error())
	}
	return logFile
}
