package setup

import (
	"gameservice/internal/di"
	"os"
)

func mustLoggerFile(cfg di.ConfigType) di.LoggerFileType {
	logFile, err := os.OpenFile(cfg.Logger.FileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		panic("Log error: " + err.Error())
	}
	return logFile
}

func deferLoggerFile(loggerFile di.LoggerFileType) {
	loggerFile.Close()
}
