package setup

import (
	"gameservice/internal/di"
)

func MustContainer(cfg di.ConfigType) di.Container {
	logger := mustLogger(cfg)
	loggerFile := mustLoggerFile(cfg)
	tracer := mustJaeger(cfg)
	db := mustDB(cfg, logger)
	cache := mustCache(cfg, logger)
	bus := mustBus(cfg, logger, cache)
	services := mustServices(db, cache, logger, bus)

	return di.Container{
		Config:     cfg,
		Logger:     logger,
		LoggerFile: loggerFile,
		Services:   services,
		DB:         db,
		Bus:        bus,
		Cache:      cache,
		Tracer:     tracer,
	}
}

func DeferContainer(container di.Container) {
	deferLoggerFile(container.LoggerFile)
	deferLogger(container.Logger)
	deferBus(container.Bus)
	deferJaeger(container.Tracer.Closer)
}
