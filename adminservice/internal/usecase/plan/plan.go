package plan

import (
	repo "adminservice/internal/adapter/db/sql"
	"adminservice/internal/di"
	"adminservice/internal/entity"
)

type UseCase struct {
	repo   repo.PlanRepo
	logger di.LoggerType
	cfg    di.ConfigType
	bus    di.Bus
}

func New(r repo.PlanRepo, logger di.LoggerType, cfg di.ConfigType, bus di.Bus) *UseCase {
	return &UseCase{
		repo:   r,
		logger: logger,
		cfg:    cfg,
		bus:    bus,
	}
}

func (u *UseCase) PlanSeason(season entity.DetailSeasonJson) error {
	dbSeason, err := season.ToDB()
	if err != nil {
		return err
	}

	// If season is not crossing with others we add it to DB
	if err := u.repo.AddNewSeason(&dbSeason); err != nil {
		return err
	}

	// Produsing new season to Core service
	go u.bus.SeasonProducer.SendSeasonInfo(dbSeason.ToDTO())

	return nil
}

func (u *UseCase) Seasons() ([]entity.Season, error) {
	seasons, err := u.repo.Seasons()
	if err != nil {
		return []entity.Season{}, err
	}

	return seasons, nil
}
