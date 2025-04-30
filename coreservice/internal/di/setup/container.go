package setup

import (
	"coreservice/internal/di"
)

func MustContainer(cfg di.ConfigType) di.Container {
	logger := mustLogger(cfg)
	loggerFile := mustLoggerFile(cfg)
	db := mustDB(cfg, logger)
	// bus := mustBus(cfg, logger)
	validator := mustValiadtor()
	services := mustServices(db, logger, cfg)

	return di.Container{
		Config:     cfg,
		Logger:     logger,
		LoggerFile: loggerFile,
		Services:   services,
		DB:         db,
		// Bus:        bus,
		Validator: validator,
	}
}

func DeferContainer(container di.Container) {
	deferLoggerFile(container.LoggerFile)
	deferLogger(container.Logger)
	deferDB()
	// deferBus(container.Bus)
}
