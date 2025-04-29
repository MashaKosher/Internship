package plan

import (
	repo "adminservice/internal/adapter/db/sql"
	"adminservice/internal/adapter/kafka/producers"
	"adminservice/internal/di"
	"adminservice/internal/entity"
	"log"
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

// Опять же, у тебя бизнес логика берет данные из параметров http хендлера, а что если поменяется веб-фрейм,
// и будут другие параметры, или вообще транспорт поменяется с http на grpc?
func (uc *UseCase) PlanSeasons(season entity.DetailSeasonJson) error {

	// var dbSeason entity.Season

	// Serialize Season JSON entity in DB entity
	// if err := pkg.StoreDeatailSeasonInDBEntity(&season, &dbSeason); err != nil {
	// 	return err
	// }

	dbSeason, err := season.ToDB()
	if err != nil {
		return err
	}

	// Finding if seasons are crossing
	if err := uc.repo.FindSeasonCross(&dbSeason); err != nil {
		return err
	}

	// If season is not crossing with others we add it to DB
	if err := uc.repo.AddNewSeason(&dbSeason); err != nil {
		return err
	}

	log.Println("Season IN DB: ", dbSeason)

	// Produsing new season to Core service
	go producers.SendSeasonInfo(dbSeason.ToDTO(), uc.cfg, uc.bus)

	return nil
}
