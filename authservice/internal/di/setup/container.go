package setup

import (
	"authservice/internal/di"
)

func MustContainer(cfg di.ConfigType) di.Container {
	logger := mustLogger(cfg)
	loggerFile := mustLoggerFile(cfg)
	tracer := mustJaeger(cfg)
	RSAKeys := mustRSAKeys(cfg, logger)
	db := mustDB(cfg, logger)
	bus := mustBus(cfg, logger, db, RSAKeys)
	services := mustServices(db, logger, RSAKeys, bus.SignUpProducer)
	validator := mustValiadtor()

	return di.Container{
		Config:     cfg,
		Logger:     logger,
		LoggerFile: loggerFile,
		Services:   services,
		DB:         db,
		Bus:        bus,
		RSAKeys:    RSAKeys,
		Validator:  validator,
		Tracer:     tracer,
	}
}

func DeferContainer(container di.Container) {
	deferLoggerFile(container.LoggerFile)
	deferLogger(container.Logger)
	deferBus(container.Bus)
	deferJaeger(container.Tracer.Closer)
}
