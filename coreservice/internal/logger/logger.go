package logger

import (
	"coreservice/internal/config"
	"os"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func CreateLogger() *os.File {
	conf := zap.NewProductionConfig()
	conf.OutputPaths = []string{"stdout", config.AppConfig.Logger.FileName}

	logger, err := conf.Build()
	if err != nil {
		panic(err)
	}

	Logger = logger
	Logger.Info("Program start")

	return creteLogFile()

}

func creteLogFile() *os.File {

	logFile, err := os.OpenFile(config.AppConfig.Logger.FileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("Log error: " + err.Error())
	}
	return logFile
}
