package setup

import (
	"adminservice/internal/di"
)

func MustContainer(cfg di.ConfigType) di.Container {
	logger := mustLogger(cfg)
	loggerFile := mustLoggerFile(cfg)
	db := mustDB(cfg, logger)
	bus := mustBus(cfg, logger)
	services := mustServices(db, logger, cfg, bus)
	validator := mustValiadtor()

	return di.Container{
		Config:     cfg,
		Logger:     logger,
		LoggerFile: loggerFile,
		Services:   services,
		DB:         db,
		Bus:        bus,
		Validator:  validator,
	}
}

func DeferContainer(container di.Container) {
	deferLoggerFile(container.LoggerFile)
	deferLogger(container.Logger)
	deferBus(container.Bus)
}
