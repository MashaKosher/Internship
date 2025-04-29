package statistic

import (
	repo "adminservice/internal/adapter/db/sql"
	"adminservice/internal/di"
)

type UseCase struct {
	repo   repo.StatisticRepo
	logger di.LoggerType
	cfg    di.ConfigType
	bus    di.Bus
}

func New(r repo.StatisticRepo, logger di.LoggerType, cfg di.ConfigType, bus di.Bus) *UseCase {
	return &UseCase{
		repo:   r,
		logger: logger,
		cfg:    cfg,
		bus:    bus,
	}
}
