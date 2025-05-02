package setup

import (
	"coreservice/internal/di"
)

func MustContainer(cfg di.ConfigType) di.Container {
	logger := mustLogger(cfg)
	loggerFile := mustLoggerFile(cfg)
	db := mustDB(cfg, logger)
	elastic := mustElastic(logger, db)
	bus := mustBus(cfg, logger, db, elastic)
	validator := mustValiadtor()
	cache := mustCache(cfg, logger)
	services := mustServices(db, logger, elastic, cache)

	return di.Container{
		Config:     cfg,
		Logger:     logger,
		LoggerFile: loggerFile,
		Services:   services,
		DB:         db,
		Bus:        bus,
		Validator:  validator,
		Elastic:    elastic,
		Cache:      cache,
	}
}

func DeferContainer(container di.Container) {
	deferLoggerFile(container.LoggerFile)
	deferLogger(container.Logger)
	deferDB()
	deferBus(container.Bus)
}
